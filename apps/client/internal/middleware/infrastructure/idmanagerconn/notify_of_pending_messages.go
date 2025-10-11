package idmanagerconn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type messageTargetRequest struct {
	TargetID  string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type notifyOfPendingMessagesRequest struct {
	MessageTargets []messageTargetRequest `json:"message_targets"`
}

func (c *idManagerConnectionImpl) NotifyOfPendingMessages(messageTargets []middleware_entities.MessageTarget) error {
	url := fmt.Sprintf("%s/pending-messages", c.baseURL)

	requestTargets := make([]messageTargetRequest, len(messageTargets))
	for i, mt := range messageTargets {
		requestTargets[i] = messageTargetRequest{
			TargetID:  string(mt.Target),
			MessageID: mt.ID,
		}
	}

	requestBody := notifyOfPendingMessagesRequest{
		MessageTargets: requestTargets,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Client-ID", string(c.ownID))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to notify of pending messages: %s", resp.Status)
	}

	return nil
}
