package idmanagerconn

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type idManagerConnectionImpl struct {
	client  *http.Client
	baseURL string
	ownID   entities.UserID
}

var _ middleware_entities.IDManagerConnection = (*idManagerConnectionImpl)(nil)

// BuildIDManagerConnection is a constructor for IDManagerConnection.
func BuildIDManagerConnection(idManagerData middleware_entities.IDManagerData, privateKey rsa.PrivateKey, ownID entities.UserID) (middleware_entities.IDManagerConnection, error) {
	//  TODO: Use private key for authentication. For now, it's ignored as per requirements.

	baseURL := fmt.Sprintf("http://%s:%d", idManagerData.IP.String(), idManagerData.Port)
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
		ownID:   ownID,
	}, nil
}
