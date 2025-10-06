package domain

import "time"

// Used to represent a client inside EagleChat system
type User struct {
	Username     string    `json:"username"`     // user readable alias
	PublicKeyPEM string    `json:"public_key"`   // PEM codified RSA Public Key
	IP           string    `json:"ip"`           // last known IP
	Connected    bool      `json:"connected"`    // current state
	LastSeen     time.Time `json:"last_seen"`    // last moment of connection/disconnection
	PendingMsgs  []string  `json:"pending_msgs"` // pending messages IDs(hashes)
}
