package p2pconnpool

import "sync"

// Close orchestrates the graceful shutdown of the entire connection pool.
func (p *p2pConnPoolImpl) Close() {
	// Use a sync.Once to prevent a panic from closing the quit channel multiple times.
	var once sync.Once
	once.Do(func() {
		close(p.quit)
	})

	// Stop the listener from accepting new connections.
	p.listener.Close()

	// Close all cached outgoing connections.
	p.mu.Lock()
	for _, conn := range p.conns {
		conn.Close()
	}
	p.mu.Unlock()
}

// Done returns a channel that is closed when the pool has fully terminated.
func (p *p2pConnPoolImpl) Done() <-chan struct{} {
	return p.done
}
