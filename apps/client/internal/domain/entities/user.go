package entities

import "eaglechat/apps/client/internal/utils/simplecrypto/rsa"

type User struct {
	ID        UserID
	Name      string
	PublicKey rsa.PublicKey
}

type UserID string

func NewUser(ID, name string, publicKey rsa.PublicKey) User {
	return User{
		ID:        UserID(ID),
		Name:      name,
		PublicKey: publicKey,
	}
}
