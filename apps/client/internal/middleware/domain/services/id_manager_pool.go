package services

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

// IDManagerPool manages listening to new id managers broadcasts in the network
type IDManagerPool interface {
	GetAny() (middleware_entities.IDManagerConnection, error)
	GetAll() ([]middleware_entities.IDManagerConnection, error)
}

type IDManagerPoolBuilder func(privateKey rsa.PrivateKey, connectionBuilder middleware_entities.IDManagerConnBuilder, ownID entities.UserID) (IDManagerPool, error)
