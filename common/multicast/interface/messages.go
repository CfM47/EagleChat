package multicast

import (
	"errors"
)

type BroadcastType string

const (
	AnnounceIDManager BroadcastType = "ANNOUNCE"
	UpdateIDManager   BroadcastType = "UPDATE"
	RegisterClient    BroadcastType = "REGISTER"
)

// BroadcastMessage defines the structure of the announcement message sent via multicast.
type BroadcastMessage struct {
	Type    BroadcastType `json:"type"`
	Time    int64         `json:"time"`
	Content []byte        `json:"content"`
}

var ErrWrongMessageType = errors.New("wrong message type")
