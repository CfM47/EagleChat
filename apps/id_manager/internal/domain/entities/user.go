package entities

import (
	"net"
	"time"
)

// IPExpirationDuration is the duration after which an IP address is considered expired
const IPExpirationDuration = 30 * time.Second

// User represents a client inside EagleChat system
type User struct {
	ID           string    `json:"id"`         // user unique id
	Username     string    `json:"username"`   // user readable alias
	PublicKeyPEM []byte    `json:"public_key"` // PEM codified RSA Public Key
	IP           *net.IP   `json:"ip"`
	LastSeen     time.Time `json:"last_seen"` // last moment of connection/disconnection
}

// NewUser creates a new User object.
// The LastSeen timestamp is rounded to zero to strip the monotonic clock reading,
// ensuring that time comparisons are deterministic and testable, especially after
// serialization and deserialization (e.g., to/from JSON), which also removes
// the monotonic clock reading.
func NewUser(id, username string, publicKeyPEM []byte, lastSeen time.Time) *User {
	return &User{
		ID:           id,
		Username:     username,
		PublicKeyPEM: publicKeyPEM,
		LastSeen:     lastSeen.UTC().Round(0),
	}
}
