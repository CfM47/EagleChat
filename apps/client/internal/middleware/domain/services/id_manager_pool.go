package services

import (
	"context"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/simplecrypto/rsa"
)

type IDManagerPool interface {
	QueryUsers(ctx context.Context, IDs []entities.UserID, omitDisconnected bool) (map[entities.UserID]middleware_entities.UserData, error)
	NotifyOfPendingMessages(context.Context, []middleware_entities.MessageTarget) error
	GetPendingMessages(context.Context) ([]middleware_entities.PendingMessage, error)
	GetRandomConnectedUsers(ctx context.Context, count int) ([]middleware_entities.UserData, error)

	Close() error
	Done() <-chan struct{}
}

type IDManagerPoolBuilder interface {
	Build(ctx context.Context, ownProfile entities.OwnProfile, CAPubkey *rsa.PublicKey) (IDManagerPool, error)
}
