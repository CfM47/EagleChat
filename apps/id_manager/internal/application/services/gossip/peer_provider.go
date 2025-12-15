package gossip

import (
	"context"
	"net"
)

// PeerProvider defines the interface for discovering service IPs.
// It aligns with common/ns.Discovery.
type PeerProvider interface {
	DiscoverIDManagerIPs(ctx context.Context) ([]net.IP, error)
}
