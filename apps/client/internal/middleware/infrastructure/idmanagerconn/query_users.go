package idmanagerconn

import (
	"bytes"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezcrypto/rsa"
	"eaglechat/common/ezlog"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type queryUsersRequest struct {
	IDs []string `json:"Ids"`
}

type userDataResponse struct {
	Name      string `json:"username"`
	PublicKey []byte `json:"public_key"`
	IP        string `json:"ip,omitempty"`
}

type queryUsersResponse map[string]userDataResponse

func (c *idManagerConnectionImpl) QueryUsers(userIDs []entities.UserID, omitDisconnected bool) (map[entities.UserID]middleware_entities.UserData, error) {
	ctx := ezlog.NewLoggerContext("query-users-request")
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
	q := req.URL.Query()
	q.Add("omit_disconnected", boolStr(omitDisconnected))
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ezlog.Log(ctx).Errorf("Failed to query users: %s", resp.Status)
		return nil, fmt.Errorf("failed to query users: %s", resp.Status)
	}

	var response queryUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		ezlog.Log(ctx).Errorf("Failed to decode query users response: %v", err)
		return nil, err
	}

	result := make(map[entities.UserID]middleware_entities.UserData)
	for id, data := range response {
		userID := entities.UserID(id)
		publicKey, err := rsa.PublicKeyFromBytes(data.PublicKey)
		if err != nil {
			ezlog.Log(ctx).Warnf("Invalid public key found while querying user '%s' from ID manager", userID)
			continue
		}

		ip := net.ParseIP(data.IP)
		if ip == nil {
			ezlog.Log(ctx).Warnf("Invalid IP address found while querying user '%s' from ID manager: %s", userID, data.IP)
			return nil, fmt.Errorf("invalid IP address for user %s: %s", userID, data.IP)
		}

		result[userID] = middleware_entities.UserData{
			User: entities.NewUser(string(userID), data.Name, *publicKey),
			IP:   &ip,
		}
	}

	return result, nil
}

func boolStr(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}
