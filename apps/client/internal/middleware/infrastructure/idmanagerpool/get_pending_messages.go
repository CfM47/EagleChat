package idmanagerpool

import (
	"context"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/lib"
)

// GetPendingMessages implements services.IDManagerPool.
func (p *idManagerPoolImpl) GetPendingMessages(ctx context.Context) ([]middleware_entities.PendingMessage, error) {
	answ := broadcast(p, ctx, getPendingMessagesAction)

	messageDict := make(map[middleware_entities.MessageTarget]middleware_entities.PendingMessage)

	for _, messagesAnsw := range answ {
		for _, message := range *messagesAnsw {
			messageDict[message.Target] = message
		}
	}

	return lib.Values(messageDict), nil
}

func getPendingMessagesAction(ctx context.Context, connection managerpool_entities.IDManagerConnection) *[]middleware_entities.PendingMessage {
	messages, err := connection.GetPendingMessages(ctx)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to get pending messages from ID manager at %s: %v", connection.BaseURL(), err)
		return nil
	}

	return &messages
}
