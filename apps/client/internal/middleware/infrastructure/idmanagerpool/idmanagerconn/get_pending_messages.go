package idmanagerconn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type messageTargetResponse struct {
	UserID    string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type getPendingMessagesResponse struct {
	MessageTargets []messageTargetResponse `json:"message_targets"`
}

func (c *idManagerConnectionImpl) GetPendingMessages(ctx context.Context) ([]middleware_entities.PendingMessage, error) {
	ezlog.Log(ctx).Info("Initiating GetPendingMessages request")

	url := fmt.Sprintf("%s/pending-messages", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to create GetPendingMessages request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		ezlog.Log(ctx).Errorf("GetPendingMessages request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ezlog.Log(ctx).Errorf("GetPendingMessages request returned non-OK status: %s", resp.Status)
		return nil, fmt.Errorf("failed to get pending messages: %s", resp.Status)
	}

	var response getPendingMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		ezlog.Log(ctx).Errorf("Failed to decode GetPendingMessages response: %v", err)
		return nil, err
	}

	result := make([]middleware_entities.PendingMessage, len(response.MessageTargets))
	for i, mt := range response.MessageTargets {
		result[i] = middleware_entities.PendingMessage{
			Target: middleware_entities.NewMessageTarget(mt.MessageID, entities.UserID(mt.UserID)),
		}
	}

	ezlog.Log(ctx).Infof("Successfully retrieved %d pending messages", len(result))

	return result, nil
}
