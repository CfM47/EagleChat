package services

import (
	"context"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/simplecrypto/rsa"
)

type Middleware interface {
	Message(ctx context.Context, target entities.User, message entities.Message) error

	QueryUser(userID entities.UserID) (entities.User, error)

	Now() time.Time
}

type Connector interface {
	Connect(ctx context.Context, listenPort uint16, ownProfile entities.OwnProfile, CAPubkey *rsa.PublicKey) (Middleware, <-chan entities.Message, error)
}

type Registerer interface {
	Register(ctx context.Context, username string, sk rsa.PrivateKey, CAPubkey *rsa.PublicKey) (entities.User, error)
}
