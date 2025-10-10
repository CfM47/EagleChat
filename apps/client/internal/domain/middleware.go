package domain

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

type Middleware interface {
	Message(userID entities.UserID, message string) error

	QueryUser(userID entities.UserID) (entities.User, error)
}

type Connector interface {
	Connect(listenPort uint16, userID entities.User, sk rsa.PrivateKey) (Middleware, <-chan entities.Message, error)
}
