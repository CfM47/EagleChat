package clientconnpool

import (
	"eaglechat/apps/client/internal/middleware/domain/entities"
)

type messageRequest struct {
	Messages []entities.PendingMessage
	Nonce    string
}

func newMessageRequest(messages []entities.PendingMessage, nonce string) messageRequest {
	return messageRequest{
		Messages: messages,
		Nonce:    nonce,
	}
}
