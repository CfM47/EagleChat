package entities

import (
	"net"
	"time"
)

// User represents a client inside EagleChat system
type User struct {
	ID           string    `json:"id"`         // user unique id
	Username     string    `json:"username"`   // user readable alias
	PublicKeyPEM string    `json:"public_key"` // PEM codified RSA Public Key
	IP           *net.IP   `json:"ip"`
	LastSeen     time.Time `json:"last_seen"` // last moment of connection/disconnection
}

// NewUser creates a new User object.
// The LastSeen timestamp is rounded to zero to strip the monotonic clock reading,
// ensuring that time comparisons are deterministic and testable, especially after
// serialization and deserialization (e.g., to/from JSON), which also removes
// the monotonic clock reading.
func NewUser(id, username, publicKeyPEM string) *User {
	return &User{
		ID:           id,
		Username:     username,
		PublicKeyPEM: publicKeyPEM,
		LastSeen:     time.Now().Round(0),
	}
}
