package jsonmessagecache

import (
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// StoreExpiring will store a message that will be immune to deletion for the
// specified `immunityPeriod`. If a message with the same ID already exists
// and is permanently immune, this operation will do nothing. If the message
// exists and is already temporarily immune, its immunity timer will be reset
// to the new `immunityPeriod`.
func (c *jsonMessageCache) StoreExpiring(message middleware_entities.PendingMessage, immunityPeriod time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, err := c.readCache()
	if err != nil {
		return err
	}

	// If a permanently immune message with the same ID exists, do nothing.
	if existing, ok := cache[message.Target]; ok && existing.IsImmune {
		return nil
	}

	cache[message.Target] = cachedMessage{
		Message:        message,
		IsImmune:       false,
		GracePeriodEnd: time.Now().Add(immunityPeriod),
	}

	return c.writeCache(cache)
}
