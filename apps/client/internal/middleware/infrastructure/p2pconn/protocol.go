package p2pconn

// MessageType is a custom type to create an enum for our wire message types.
type MessageType string

const (
	// MsgTypeData indicates a message containing application data.
	MsgTypeData MessageType = "data"
	// MsgTypeAck indicates a message acknowledging receipt of data.
	MsgTypeAck MessageType = "ack"
)

// WireMessage is the low-level object that wraps all P2P communication.
// It is marshalled to JSON for transport.
type WireMessage struct {
	Type    MessageType `json:"type"`
	ID      string      `json:"id"`
	Payload []byte      `json:"payload,omitempty"`
}
