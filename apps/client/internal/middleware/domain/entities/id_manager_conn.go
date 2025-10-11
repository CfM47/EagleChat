package entities

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

type IDManagerConnection interface {
	QueryUsers([]entities.UserID) (map[entities.UserID]UserData, error)
	NotifyOfPendingMessages([]MessageTarget) error
	GetPendingMessages() ([]PendingMessage, error)
}

type IDManagerConnBuilder func(IDManagerData, rsa.PrivateKey, entities.UserID) (IDManagerConnection, error)
