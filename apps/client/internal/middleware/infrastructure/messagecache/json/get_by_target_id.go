package jsonmessagecache

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// GetByTargetId returns all pending messages for a specific recipient.
func (c *jsonMessageCache) GetByTargetId(userID entities.UserID) []middleware_entities.PendingMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cache, err := c.readCache()
	if err != nil {
		return []middleware_entities.PendingMessage{}
	}

	var results []middleware_entities.PendingMessage
	for _, cachedMsg := range cache {
		if cachedMsg.Message.Target.TargetID == userID {
			results = append(results, cachedMsg.Message)
		}
	}

	return results
}
