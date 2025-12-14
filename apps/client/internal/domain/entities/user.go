package entities

import (
	"time"

	"eaglechat/common/simplecrypto/rsa"
)

type User struct {
	ID        UserID
	Name      string
	PublicKey rsa.PublicKey
	LastSeen  time.Time
}

type UserID string

func NewUser(ID, name string, publicKey rsa.PublicKey, lastSeen time.Time) User {
	return User{
		ID:        UserID(ID),
		Name:      name,
		PublicKey: publicKey,
		LastSeen:  lastSeen,
	}
}
