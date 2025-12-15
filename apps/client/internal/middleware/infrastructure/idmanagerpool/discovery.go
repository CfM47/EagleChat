package idmanagerpool

import (
	"context"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerverifier"
	"eaglechat/common/ezlog"
	"eaglechat/common/ns"
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
		pk, _, err := idmanagerverifier.VerifyIDManager(ctx, ip, middleware_entities.DefaultIDManagerPort, p.CAPubkey)
		if err != nil {
			ezlog.Log(ctx).Warnf("ID manager verification for DNS looked up id manager failed at ip %s", ip)
		} else {
			p.repository.Add(middleware_entities.NewIDManagerData(ip, middleware_entities.DefaultIDManagerPort, *pk))
		}
	}
}
