package ns

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
)

// cachingDiscovery is a decorator for ns.Discovery that adds a file-based cache.
type cachingDiscovery struct {
	wrappedDiscovery Discovery
	cacheFilePath    string
}

// NewCachingDiscovery creates a new cachingDiscovery.
func NewCachingDiscovery(wrappedDiscovery Discovery, cacheFilePath string) Discovery {
	return &cachingDiscovery{
		wrappedDiscovery: wrappedDiscovery,
		cacheFilePath:    cacheFilePath,
	}
}

// GetIPs attempts to get IPs from the wrapped discovery service.
// If it fails, it falls back to the cache.
// If it succeeds, it updates the cache.
func (d *cachingDiscovery) DiscoverIDManagerIPs(ctx context.Context) ([]net.IP, error) {
	ips, err := d.wrappedDiscovery.DiscoverIDManagerIPs(ctx)
	if err == nil {
		d.updateCache(ips)
		return ips, nil
	}

	// If wrapped discovery fails, try to read from cache
	cachedIPs, cacheErr := d.readCache()
	if cacheErr != nil {
		// If cache also fails, return the original error from wrappedDiscovery
		return nil, err
	}
	return cachedIPs, nil
}

func (d *cachingDiscovery) GetOwnHost() (string, error) {
	return d.wrappedDiscovery.GetOwnHost()
}

func (d *cachingDiscovery) updateCache(ips []net.IP) {
	data, err := json.Marshal(ips)
	if err != nil {
		// Log the error, but don't fail the operation
		return
	}
	// Ensure the directory exists
	dir := filepath.Dir(d.cacheFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// Log the error
		return
	}
	os.WriteFile(d.cacheFilePath, data, 0644)
}

func (d *cachingDiscovery) readCache() ([]net.IP, error) {
	data, err := os.ReadFile(d.cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []net.IP{}, nil // Return empty slice if file doesn't exist, not an error
		}
		return nil, err
	}

	var ips []net.IP
	if err := json.Unmarshal(data, &ips); err != nil {
		return nil, err
	}
	return ips, nil
}
