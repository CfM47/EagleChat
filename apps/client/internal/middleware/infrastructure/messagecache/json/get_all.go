package jsonmessagecache

import (
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// GetAll returns all messages currently in the cache.
func (c *jsonMessageCache) GetAll() []middleware_entities.PendingMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cache, err := c.readCache()
	if err != nil {
		return []middleware_entities.PendingMessage{}
	}

	messages := make([]middleware_entities.PendingMessage, 0, len(cache))
	for _, cachedMsg := range cache {
		messages = append(messages, cachedMsg.Message)
	}

	return messages
}
