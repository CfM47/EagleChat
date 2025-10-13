package p2pconn

import (
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"fmt"
	"net"
	"sync"
)

type p2pListener struct {
	listener    net.Listener
	connections chan entities.P2PConnection
	quit        chan struct{}
	wg          sync.WaitGroup
}

// StartListener creates and starts a new P2P listener on the given port.
func StartListener(port uint16) (entities.P2PConnListener, error) {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	l := &p2pListener{
		listener:    ln,
		connections: make(chan entities.P2PConnection),
		quit:        make(chan struct{}),
	}

	l.wg.Add(1)
	go l.acceptLoop()

	return l, nil
}

func (l *p2pListener) acceptLoop() {
	defer l.wg.Done()
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			select {
			case <-l.quit:
				return // Graceful shutdown
			default:
				// FIXME: Handle unexpected error
				return
			}
		}

		// Wrap the new connection and send it to the application
		p2pConn := newP2PConnection(conn)
		select {
		case l.connections <- p2pConn:
		case <-l.quit:
			// If a shutdown is initiated while waiting to send, close the new conn.
			conn.Close()
			return
		}
	}
}

func (l *p2pListener) Connections() <-chan entities.P2PConnection {
	return l.connections
}

func (l *p2pListener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *p2pListener) Close() error {
	close(l.quit)
	err := l.listener.Close()
	l.wg.Wait() // Wait for the accept loop to finish
	close(l.connections)
	return err
}
