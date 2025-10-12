package entities

import "eaglechat/apps/client/internal/utils/simplecrypto/rsa"

// OwnProfile represents the client's own user identity, including the private key.
type OwnProfile struct {
	User       User
	PrivateKey rsa.PrivateKey
}

func NewOwnProfile(user User, sk rsa.PrivateKey) OwnProfile {
	return OwnProfile{
		User:       user,
		PrivateKey: sk,
	}
}
