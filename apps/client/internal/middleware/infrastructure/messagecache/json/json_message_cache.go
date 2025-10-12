package jsonmessagecache

import (
	"eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// jsonMessageCache implements the MessageCache interface using a JSON file.
type jsonMessageCache struct {
	path string
	mu   sync.RWMutex
}

// cachedMessage is the internal representation of a message in the cache.
type cachedMessage struct {
	Message        middleware_entities.PendingMessage `json:"message"`
	IsImmune       bool                               `json:"is_immune"`
	GracePeriodEnd time.Time                          `json:"grace_period_end"`
}

// NewJSONMessageCache creates a new file-based message cache.
func NewJSONMessageCache(path string) (messagecache.MessageCache, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &jsonMessageCache{
		path: path,
	}, nil
}

// readCache safely reads and deserializes the cache from disk.
func (c *jsonMessageCache) readCache() (map[middleware_entities.MessageTarget]cachedMessage, error) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[middleware_entities.MessageTarget]cachedMessage), nil
		}
		return nil, err
	}

	var cache map[middleware_entities.MessageTarget]cachedMessage
	if err := json.Unmarshal(data, &cache); err != nil {
		return make(map[middleware_entities.MessageTarget]cachedMessage), nil // Return empty map on corruption
	}
	return cache, nil
}

// writeCache safely serializes and writes the cache to disk.
func (c *jsonMessageCache) writeCache(cache map[middleware_entities.MessageTarget]cachedMessage) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}
