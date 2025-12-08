package idmanagerconn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
)

type idManagerConnectorImpl struct{}

func NewIDManagerConnector() managerpool_entities.IDManagerConnector {
	return &idManagerConnectorImpl{}
}

var _ managerpool_entities.IDManagerConnector = (*idManagerConnectorImpl)(nil)

// Connect implements entities.IDManagerConnector.
func (i *idManagerConnectorImpl) Connect(ctx context.Context, ownProfile entities.OwnProfile, IDManagerData middleware_entities.IDManagerData) (managerpool_entities.IDManagerConnection, error) {
	//  TODO: Use private key for authentication. For now, it's ignored.
	//  FIXME: add logging

	baseURL := fmt.Sprintf("http://%s:%d", IDManagerData.IP.String(), IDManagerData.Port)
	client := &http.Client{}

	// Health check
	resp, err := client.Get(baseURL + "/status")
	if err != nil {
		return nil, fmt.Errorf("id manager connection failed health check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("id manager connection failed health check: status code %d", resp.StatusCode)
	}

	var statusResponse struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		return nil, fmt.Errorf("id manager connection failed health check: could not decode response: %w", err)
	}

	if statusResponse.Status != "ok" {
		return nil, fmt.Errorf("id manager connection failed health check: unexpected status '%s'", statusResponse.Status)
	}

	return &idManagerConnectionImpl{
		client:  client,
		baseURL: baseURL,
		ownID:   ownProfile.User.ID,
	}, nil
}
