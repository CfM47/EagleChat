package jsonmessagecache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
	"eaglechat/common/lib"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// jsonMessageCache implements the MessageCache interface using a JSON file.
type jsonMessageCache struct {
	path string
	mu   sync.RWMutex
}

var _ messagecache.MessageCache = (*jsonMessageCache)(nil)

// cachedMessage is the internal representation of a message in the cache.
type cachedMessage struct {
	Message        middleware_entities.PendingMessage `json:"message"`
	IsImmune       bool                               `json:"is_immune"`
	GracePeriodEnd time.Time                          `json:"grace_period_end"`
}

// NewJSONMessageCache creates a new file-based message cache.
func NewJSONMessageCache(path string) (messagecache.MessageCache, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
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

	var mappedCache map[string]cachedMessage
	if err := json.Unmarshal(data, &mappedCache); err != nil {
		return make(map[middleware_entities.MessageTarget]cachedMessage), nil // Return empty map on corruption
	}

	cache := lib.MapKeys(mappedCache, func(ts string) middleware_entities.MessageTarget {
		parts := strings.Split(ts, ":")
		if len(parts) != 2 {
			panic("invalid message cache")
		}

		return middleware_entities.NewMessageTarget(parts[1], entities.UserID(parts[0]))
	})

	return cache, nil
}

// writeCache safely serializes and writes the cache to disk.
func (c *jsonMessageCache) writeCache(cache map[middleware_entities.MessageTarget]cachedMessage) error {
	mappedCache := lib.MapKeys(cache, func(t middleware_entities.MessageTarget) string {
		return string(t.TargetID) + ":" + t.MessageID
	})

	data, err := json.MarshalIndent(mappedCache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}
