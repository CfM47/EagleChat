package jsonmessagecache

import (
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/lib"
)

func (c *jsonMessageCache) GetByTargets(targets []middleware_entities.MessageTarget) []middleware_entities.PendingMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	targetSet := lib.MakeSet(targets)

	cache, err := c.readCache()
	if err != nil {
		return []middleware_entities.PendingMessage{}
	}

	var results []middleware_entities.PendingMessage
	for _, cachedMsg := range cache {
		if _, ok := targetSet[cachedMsg.Message.Target]; ok {
			results = append(results, cachedMsg.Message)
		}
	}

	return results
}
