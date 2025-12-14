package idmanagerconn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type queryUsersRequest struct {
	IDs []string `json:"Ids"`
}

type queryUsersResponse map[string]userData

func (c *idManagerConnectionImpl) QueryUsers(ctx context.Context, userIDs []entities.UserID, omitDisconnected bool) (map[entities.UserID]middleware_entities.UserData, error) {
	url := fmt.Sprintf("%s/users", c.baseURL)
	ezlog.Log(ctx).Infof("Querying users to %s", url)

	stringIDs := make([]string, len(userIDs))
	for i, id := range userIDs {
		stringIDs[i] = string(id)
	}

	requestBody := queryUsersRequest{
		IDs: stringIDs,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to marshal request: %v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to create request: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("omit_disconnected", boolStr(omitDisconnected))
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to perform request: %v", err)
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
	for _, data := range response {
		userData, err := buildUserData(ctx, data)
		if err != nil {
			return nil, err
		}

		result[userData.ID] = userData
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
