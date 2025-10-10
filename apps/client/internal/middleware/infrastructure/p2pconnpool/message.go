package p2pconnpool

import (
	"fmt"
	"net"
	"strconv"
)

// Message sends content to a peer, creating and caching the connection if necessary.
func (p *p2pConnPoolImpl) Message(ip, portStr string, content []byte) error {
	addrKey := net.JoinHostPort(ip, portStr)

	// Fast path: Check if connection exists with a read lock.
	p.mu.RLock()
	conn, ok := p.conns[addrKey]
	p.mu.RUnlock()

	if ok {
		return conn.Send(content)
	}

	// Slow path: Connection doesn't exist, so get a write lock to create it.
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check, in case another goroutine created the connection while we waited for the lock.
	conn, ok = p.conns[addrKey]
	if ok {
		return conn.Send(content)
	}

	// Parse port string to uint16 for the dialer.
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port: %s", portStr)
	}

	// Dial the new connection.
	newConn, err := p.dialer(ip, uint16(port))
	if err != nil {
		return err
	}

	// Add the new connection to the cache.
	p.conns[addrKey] = newConn

	// Start a forwarder to manage the new connection.
	p.wg.Add(1)
	go p.forwarder(newConn, addrKey)

	return newConn.Send(content)
}
