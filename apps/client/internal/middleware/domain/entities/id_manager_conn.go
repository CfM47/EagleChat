package entities

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezcrypto/rsa"
)

type IDManagerConnection interface {
	QueryUsers(IDs []entities.UserID, omitDisconnected bool) (map[entities.UserID]UserData, error)
	NotifyOfPendingMessages([]MessageTarget) error
	GetPendingMessages() ([]PendingMessage, error)
}

type IDManagerConnBuilder func(data IDManagerData, ownSk rsa.PrivateKey, ownID entities.UserID) (IDManagerConnection, error)
