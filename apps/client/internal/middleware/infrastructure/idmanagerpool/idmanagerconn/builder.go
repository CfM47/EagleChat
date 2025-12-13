package idmanagerconn

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

const checkHealthTimeout = 5 * time.Second

type idManagerConnectorImpl struct {
	ownProfile entities.OwnProfile
}

func NewIDManagerConnector(ownProfile entities.OwnProfile, CAPubkey rsa.PublicKey) managerpool_entities.IDManagerConnector {
	return &idManagerConnectorImpl{
		ownProfile: ownProfile,
	}
}

var _ managerpool_entities.IDManagerConnector = (*idManagerConnectorImpl)(nil)

// Connect implements entities.IDManagerConnector.
func (c *idManagerConnectorImpl) Connect(ctx context.Context, IDManagerData middleware_entities.IDManagerData) (managerpool_entities.IDManagerConnection, error) {
	baseURL := fmt.Sprintf("http://%s:%d", IDManagerData.IP.String(), IDManagerData.Port)
	client, ok := checkHealth(ctx, IDManagerData.IP, IDManagerData.Port)

	if !ok {
		return nil, fmt.Errorf("ID manager at %s failed health check", baseURL)
	}

	return NewIDManagerConnectionImpl(client, baseURL, c.ownProfile.User.ID, IDManagerData.PublicKey), nil
}

type StatusResponse struct {
	Status string `json:"status"`
}

func checkHealth(ctx context.Context, ip net.IP, port uint16) (*http.Client, bool) {
	url := fmt.Sprintf("http://%s:%d/status", ip.String(), port)
	ezlog.Log(ctx).Infof("Fetching status from %s", url)

	reqCtx, cancel := context.WithTimeout(ctx, checkHealthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		ezlog.Log(ctx).Warnf("Error creating request for %s: %v", url, err)
		return nil, false
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		ezlog.Log(ctx).Warnf("Error fetching status from %s: %v", url, err)
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ezlog.Log(ctx).Warnf("ID Manager at %s returned non-200 status: %s", url, resp.Status)
		return nil, false
	}

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		ezlog.Log(ctx).Warnf("Error decoding metadata from %s: %v", url, err)
		return nil, false
	}
	if status.Status != "ok" {
		ezlog.Log(ctx).Warnf("ID Manager at %s returned non-ok status: %s", url, status.Status)
		return nil, false
	}

	return &client, true
}
