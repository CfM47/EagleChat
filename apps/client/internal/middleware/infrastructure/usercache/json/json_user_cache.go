package jsonusercache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
)

// jsonUserCache implements the UserCacheRepository interface using a JSON file.
type jsonUserCache struct {
	path           string
	expirationTime time.Duration
	mu             sync.RWMutex
}

// NewJSONUserCache creates a new file-based user cache.
func NewJSONUserCache(path string, expirationTime time.Duration) (usercache.UserCacheRepository, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &jsonUserCache{
		path:           path,
		expirationTime: expirationTime,
	}, nil
}

// readCache safely reads and deserializes the cache from disk.
func (c *jsonUserCache) readCache() (map[string]middleware_entities.UserData, error) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]middleware_entities.UserData), nil
		}
		return nil, err
	}

	var cache map[string]middleware_entities.UserData
	if err := json.Unmarshal(data, &cache); err != nil {
		// On corruption, return an empty map to avoid losing all data
		return make(map[string]middleware_entities.UserData), nil
	}
	return cache, nil
}

// writeCache safely serializes and writes the cache to disk.
func (c *jsonUserCache) writeCache(cache map[string]middleware_entities.UserData) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}
