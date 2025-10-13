package p2pconnpool

import (
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"sync"
)

// p2pConnPoolImpl is a thread-safe implementation of the P2PConnPool service.
type p2pConnPoolImpl struct {
	dialer   entities.P2PDialer
	listener entities.P2PConnListener

	conns map[string]entities.P2PConnection
	mu    sync.RWMutex

	incoming chan []byte
	quit     chan struct{}
	done     chan struct{}
	wg       sync.WaitGroup
}

// BuildP2PConnPool creates and initializes a new P2P connection pool.
func BuildP2PConnPool(dialer entities.P2PDialer, listenerStarter entities.P2PListenStarter, listenPort uint16) (services.P2PConnPool, error) {
	listener, err := listenerStarter(listenPort)
	if err != nil {
		return nil, err
	}

	p := &p2pConnPoolImpl{
		dialer:   dialer,
		listener: listener,
		conns:    make(map[string]entities.P2PConnection),
		incoming: make(chan []byte, 128), // Buffered channel
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	p.wg.Add(1) // For the acceptLoop
	go p.acceptLoop()

	// This goroutine waits for all other goroutines to finish before closing the done channel.
	go func() {
		p.wg.Wait()
		close(p.done)
	}()

	return p, nil
}

// acceptLoop ranges over the listener's connections channel and starts a forwarder for each.
func (p *p2pConnPoolImpl) acceptLoop() {
	defer p.wg.Done()
	for {
		select {
		case conn, ok := <-p.listener.Connections():
			if !ok {
				return // Channel closed
			}
			p.wg.Add(1)              // Add to waitgroup for the new forwarder
			go p.forwarder(conn, "") // No address key for incoming conns
		case <-p.quit:
			return
		}
	}
}

// forwarder manages a single connection, funneling its messages to the pool's
// incoming channel and handling its termination.
func (p *p2pConnPoolImpl) forwarder(conn entities.P2PConnection, addrKey string) {
	defer p.wg.Done()
	for {
		select {
		case msg, ok := <-conn.Receive():
			if !ok {
				return // Connection's receive channel closed
			}
			// Forward message to the pool's central channel
			select {
			case p.incoming <- msg:
			case <-p.quit:
				return
			}
		case <-conn.Done():
			// The connection closed on its own.
			if addrKey != "" { // Only outgoing connections are in the map
				p.mu.Lock()
				delete(p.conns, addrKey)
				p.mu.Unlock()
			}
			return
		case <-p.quit:
			// The entire pool is shutting down.
			return
		}
	}
}
