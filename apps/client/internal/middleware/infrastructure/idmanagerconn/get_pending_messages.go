package idmanagerconn

import (
	"eaglechat/apps/client/internal/domain/entities"
	"encoding/json"
	"fmt"
	"net/http"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type messageTargetResponse struct {
	UserID    string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type getPendingMessagesResponse struct {
	MessageTargets []messageTargetResponse `json:"message_targets"`
}

func (c *idManagerConnectionImpl) GetPendingMessages() ([]middleware_entities.PendingMessage, error) {
	url := fmt.Sprintf("%s/pending-messages", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, fmt.Errorf("failed to get pending messages: %s", resp.Status)
	}

	var response getPendingMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	result := make([]middleware_entities.PendingMessage, len(response.MessageTargets))
	for i, mt := range response.MessageTargets {
		result[i] = middleware_entities.PendingMessage{
			Target: middleware_entities.NewMessageTarget(mt.MessageID, entities.UserID(mt.UserID)),
		}
	}

	return result, nil
}
