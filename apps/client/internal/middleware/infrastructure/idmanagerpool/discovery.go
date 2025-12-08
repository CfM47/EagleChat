package idmanagerpool

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/ns"
)

const (
	checkHealthTimeout = 5 * time.Second
)

func (p *idManagerPoolImpl) pollDNSLoop(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting DNS polling loop for ID Managers...")

	defer close(p.doneChan)
	defer ezlog.Log(ctx).Info("Stopped DNS polling loop for ID Managers.")

	ticker := time.NewTicker(DNSPOllInterval)
	defer ticker.Stop()

	// Poll once immediately on startup
	pollCtx := ezlog.NewLoggerContext("id-manager-pool-poll")
	p.pollDNS(pollCtx)

	for {
		select {
		case <-ticker.C:
			pollCtx = ezlog.NewLoggerContext("id-manager-pool-poll")
			p.pollDNS(pollCtx)
		case <-p.quitChan:
			return
		}
	}
}

func (p *idManagerPoolImpl) pollDNS(ctx context.Context) {
	ips, err := ns.NewDNSDiscovery().DiscoverIDManagerIPs(ctx)
	if err != nil {
		ezlog.Log(ctx).Warnf("DNS lookup for ID managers failed: %v", err)
		return
	}

	for _, ip := range ips {
		if checkHealth(ctx, ip, middleware_entities.DefaultIDManagerPort) {
			p.repository.Add(middleware_entities.NewIDManagerData(ip, middleware_entities.DefaultIDManagerPort))
		} else {
			ezlog.Log(ctx).Warnf("CheckHealth for DNS looked up id manager failed at ip %s", ip)
		}
	}
}

func checkHealth(ctx context.Context, ip net.IP, port uint16) bool {
	url := fmt.Sprintf("http://%s:%d/status", ip.String(), port)
	ezlog.Log(ctx).Infof("Fetching status from %s", url)

	reqCtx, cancel := context.WithTimeout(ctx, checkHealthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		ezlog.Log(ctx).Warnf("Error creating request for %s: %v", url, err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ezlog.Log(ctx).Warnf("Error fetching status from %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ezlog.Log(ctx).Warnf("ID Manager at %s returned non-200 status: %s", url, resp.Status)
		return false
	}

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		ezlog.Log(ctx).Warnf("Error decoding metadata from %s: %v", url, err)
		return false
	}
	if status.Status != "ok" {
		ezlog.Log(ctx).Warnf("ID Manager at %s returned non-ok status: %s", url, status.Status)
		return false
	}

	return true
}
