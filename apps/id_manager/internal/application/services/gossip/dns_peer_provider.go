package gossip

import (
	"context"
	"eaglechat/common/ns"
	"net"
)

// DnsPeerProvider implements PeerProvider using a wrapped ns.Discovery (e.g., cached or raw DNS).
type DnsPeerProvider struct {
	discovery ns.Discovery
}

// NewDnsPeerProvider creates a new DnsPeerProvider.
func NewDnsPeerProvider(discovery ns.Discovery) *DnsPeerProvider {
	return &DnsPeerProvider{discovery: discovery}
}

// DiscoverIDManagerIPs discovers ID Manager IPs using the wrapped ns.Discovery.
func (p *DnsPeerProvider) DiscoverIDManagerIPs(ctx context.Context) ([]net.IP, error) {
	return p.discovery.DiscoverIDManagerIPs(ctx)
}
