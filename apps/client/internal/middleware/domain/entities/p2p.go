package entities

import "net"

// P2PConnection represents an active, message-oriented P2P connection to a peer.
// It abstracts away the underlying network protocol and provides a simple
// interface for sending and receiving byte slices.
type P2PConnection interface {
	// Send transmits a message to the peer and blocks until delivery is
	// confirmed via an application-level acknowledgement from the remote peer.
	//
	// A nil error provides a strong guarantee that the message has been received
	// and consumed by the remote peer's connection layer. It does NOT guarantee
	// that the remote application has fully processed the message.
	//
	// An error is returned if this confirmation is not received within a
	// predefined timeout, or if the connection is closed.
	Send(data []byte) error

	// Receive returns a read-only channel that delivers incoming messages
	// from the peer.
	Receive() <-chan []byte

	// Close gracefully shuts down the connection and all associated goroutines.
	Close()

	// Done returns a channel that is closed when the connection is fully terminated.
	// This is useful for waiting on the connection to shut down.
	Done() <-chan struct{}
}

// P2PDialer defines the function signature for creating a new P2P connection.
// This allows for dependency injection of the connection logic.
type P2PDialer func(ip string, port uint16) (P2PConnection, error)

// P2PConnListener handles the process of accepting incoming peer connections.
type P2PConnListener interface {
	// Connections returns a read-only channel that provides newly accepted
	// peer connections.
	Connections() <-chan P2PConnection

	// Close stops the listener and releases the network port.
	Close() error

	// Addr returns the local network address the listener is bound to.
	Addr() net.Addr
}

// P2PListenStarter defines the function signature for creating a new P2P listener.
type P2PListenStarter func(port uint16) (P2PConnListener, error)
