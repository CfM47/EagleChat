package entities

import (
	"eaglechat/common/simplecrypto/rsa"
	"time"
)

type User struct {
	ID        UserID
	Name      string
	PublicKey rsa.PublicKey
	LastSeen  time.Time
}

type UserID string

func NewUser(ID, name string, publicKey rsa.PublicKey) User {
	return User{
		ID:        UserID(ID),
		Name:      name,
		PublicKey: publicKey,
		LastSeen:  time.Now().UTC().Round(0),
	}
}
