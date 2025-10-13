package ns

import (
	"context"
	"net"
)

// Discovery defines the interface for discovering service IPs.
type Discovery interface {
	// DiscoverIDManagerIPs returns a slice of IP addresses for available ID Manager instances.
	DiscoverIDManagerIPs(ctx context.Context) ([]net.IP, error)
}
