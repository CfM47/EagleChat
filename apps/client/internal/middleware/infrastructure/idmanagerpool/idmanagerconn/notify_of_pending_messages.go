package idmanagerconn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type messageTargetRequest struct {
	TargetID  string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type notifyOfPendingMessagesRequest struct {
	MessageTargets []messageTargetRequest `json:"message_targets"`
}

func (c *idManagerConnectionImpl) NotifyOfPendingMessages(ctx context.Context, ownID entities.UserID, messageTargets []middleware_entities.MessageTarget) error {
	//  FIXME: add logging

	url := fmt.Sprintf("%s/pending-messages", c.baseURL)

	requestTargets := make([]messageTargetRequest, len(messageTargets))
	for i, mt := range messageTargets {
		requestTargets[i] = messageTargetRequest{
			TargetID:  string(mt.TargetID),
			MessageID: mt.MessageID,
		}
	}

	requestBody := notifyOfPendingMessagesRequest{
		MessageTargets: requestTargets,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Client-ID", string(ownID))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to notify of pending messages: %s", resp.Status)
	}

	return nil
}
