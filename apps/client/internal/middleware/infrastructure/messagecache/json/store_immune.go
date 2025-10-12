package jsonmessagecache

import (
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// StoreImmune will store a message in the cache that will not be deleted until DeleteImmune
// is called with its id.
func (c *jsonMessageCache) StoreImmune(message middleware_entities.PendingMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, err := c.readCache()
	if err != nil {
		return err
	}

	cache[message.Target] = cachedMessage{
		Message:        message,
		IsImmune:       true,
		GracePeriodEnd: time.Time{}, // No grace period for immune messages
	}

	return c.writeCache(cache)
}
