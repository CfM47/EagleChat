package discovery

import (
	"context"
	"fmt"
	"net"
)

const (
	// IDManagerServiceAddress is the DNS name for the ID Manager service in Docker Swarm.
	IDManagerServiceAddress = "tasks.id_manager"
)

// dnsDiscovery implements the Discovery interface using DNS lookups.
type dnsDiscovery struct {
	serviceAddress string
}

// NewDNSDiscovery creates a new DNS-based discovery service.
func NewDNSDiscovery() Discovery {
	return &dnsDiscovery{
		serviceAddress: IDManagerServiceAddress,
	}
}

// DiscoverIDManagerIPs performs a DNS lookup to find the IP addresses of ID Manager tasks.
func (d *dnsDiscovery) DiscoverIDManagerIPs(ctx context.Context) ([]net.IP, error) {
	ips, err := net.LookupIP("id_manager")
	if err != nil {
		return nil, fmt.Errorf("dns lookup for %s failed: %w", d.serviceAddress, err)
	}
	return ips, nil
}
