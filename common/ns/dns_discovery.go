package ns

import (
	"context"
	"fmt"
	"net"
)

const (
	// IDManagerDefaultAlias is the default DNS name for the ID Manager instances.
	IDManagerDefaultAlias = "id_manager.eaglechat.local"
)

// dnsDiscovery implements the Discovery interface using DNS lookups.
type dnsDiscovery struct {
	roleAlias string
}

// NewDNSDiscovery creates a new DNS-based discovery service.
func NewDNSDiscovery() Discovery {
	return &dnsDiscovery{
		roleAlias: IDManagerDefaultAlias,
	}
}

// DiscoverIDManagerIPs performs a DNS lookup to find the IP addresses of ID Manager tasks.
func (d *dnsDiscovery) DiscoverIDManagerIPs(ctx context.Context) ([]net.IP, error) {
	ips, err := net.LookupIP(d.roleAlias)
	if err != nil {
		return nil, fmt.Errorf("dns lookup for %s failed: %w", d.roleAlias, err)
	}
	return ips, nil
}
