package entities

import (
	"context"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

type IDManagerConnection interface {
	QueryUsers(ctx context.Context, IDs []entities.UserID, omitDisconnected bool) (map[entities.UserID]middleware_entities.UserData, error)
	NotifyOfPendingMessages(ctx context.Context, ownID entities.UserID, pendingMessageTargets []middleware_entities.MessageTarget) error
	GetPendingMessages(context.Context) ([]middleware_entities.PendingMessage, error)
	GetRandomConnectedUsers(ctx context.Context, count int) ([]middleware_entities.UserData, error)

	BaseURL() string
}

type IDManagerConnector interface {
	Connect(ctx context.Context, IDManagerData middleware_entities.IDManagerData) (IDManagerConnection, error)
}
