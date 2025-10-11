package idmanagerconn

import (
	"bytes"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type queryUsersRequest struct {
	IDs []string `json:"ids"`
}

type userDataResponse struct {
	Name      string `json:"name"`
	PublicKey []byte `json:"public_key"`
	IP        string `json:"ip,omitempty"`
}

type queryUsersResponse map[string]userDataResponse

func (c *idManagerConnectionImpl) QueryUsers(userIDs []entities.UserID) (map[entities.UserID]middleware_entities.UserData, error) {
	url := fmt.Sprintf("%s/users", c.baseURL)

	stringIDs := make([]string, len(userIDs))
	for i, id := range userIDs {
		stringIDs[i] = string(id)
	}

	requestBody := queryUsersRequest{
		IDs: stringIDs,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to query users: %s", resp.Status)
	}

	var response queryUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	result := make(map[entities.UserID]middleware_entities.UserData)
	for id, data := range response {
		userID := entities.UserID(id)
		publicKey, err := rsa.PublicKeyFromBytes(data.PublicKey)
		if err != nil {
			log.Printf("invalid public key found while querying user '%s' from id manager", userID)
			continue
		}

		ip := net.ParseIP(data.IP)

		result[userID] = middleware_entities.UserData{
			User: entities.NewUser(string(userID), data.Name, *publicKey),
			IP:   &ip,
		}
	}

	return result, nil
}
