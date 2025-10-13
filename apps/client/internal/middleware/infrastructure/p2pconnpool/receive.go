package p2pconnpool

// Receive returns the channel that aggregates incoming messages from all connections.
func (p *p2pConnPoolImpl) Receive() <-chan []byte {
	return p.incoming
}
