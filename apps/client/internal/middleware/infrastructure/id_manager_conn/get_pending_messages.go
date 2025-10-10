package id_manager_conn

import (
	"bytes"
	"eaglechat/apps/client/internal/domain/entities"
	"encoding/json"
	"fmt"
	"net/http"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type getPendingMessagesRequest struct {
	TargetID *string `json:"target_id"`
}

type messageTargetResponse struct {
	TargetID  string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type getPendingMessagesResponse struct {
	MessageTargets []messageTargetResponse `json:"message_targets"`
}

func (c *idManagerConnectionImpl) GetPendingMessages() ([]middleware_entities.MessageTarget, error) {
	url := fmt.Sprintf("%s/pending-messages", c.baseURL)

	requestBody := getPendingMessagesRequest{
		TargetID: nil,
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
		return nil, fmt.Errorf("failed to get pending messages: %s", resp.Status)
	}

	var response getPendingMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	result := make([]middleware_entities.MessageTarget, len(response.MessageTargets))
	for i, mt := range response.MessageTargets {
		result[i] = middleware_entities.NewMessageTarget(mt.MessageID, entities.UserID(mt.TargetID))
	}

	return result, nil
}
