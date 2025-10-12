package services

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/simplecrypto/rsa"
)

type Middleware interface {
	Message(target entities.User, message entities.Message) error

	QueryUser(userID entities.UserID) (entities.User, error)
}

type Connector interface {
	Connect(listenPort uint16, userID entities.User, sk rsa.PrivateKey) (Middleware, <-chan entities.Message, error)
}

type Registerer interface {
	Register(ctx context.Context, username string, sk rsa.PrivateKey) (entities.User, error)
}
