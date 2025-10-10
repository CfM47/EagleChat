package services

import "eaglechat/apps/client/internal/middleware/domain/entities"

// P2PConnPool manages a cache of active P2P connections and aggregates all
// incoming messages into a single channel.
type P2PConnPool interface {
	// Message sends content to a specific peer, identified by IP and port.
	// It will reuse a cached connection or dial a new one if necessary.
	//
	// This is a blocking call that waits for confirmation of delivery from the
	// remote peer, inheriting the synchronous, ACK-based behavior of the
	// underlying P2PConnection.Send method. A nil error guarantees the message
	// has been consumed by the remote peer's connection layer.
	Message(ip, port string, content []byte) error

	// Receive returns a channel that aggregates incoming messages from all
	// active connections in the pool.
	Receive() <-chan []byte

	// Close gracefully shuts down the pool, the listener, and all underlying
	// connections.
	Close()

	// Done returns a channel that is closed only when the pool has fully
	// terminated.
	Done() <-chan struct{}
}

type P2PConnPoolBuilder func(dialer entities.P2PDialer, listenerStarter entities.P2PListenStarter, listenPort uint16) (P2PConnPool, error)
