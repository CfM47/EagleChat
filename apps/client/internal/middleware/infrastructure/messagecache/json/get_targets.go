package jsonmessagecache

import (
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
)

// GetTargets returns lists of immune and non-immune message targets.
func (c *jsonMessageCache) GetTargets() messagecache.PendingMessageTargetLists {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cache, err := c.readCache()
	if err != nil {
		return messagecache.PendingMessageTargetLists{}
	}

	lists := messagecache.PendingMessageTargetLists{
		Immune:    make([]middleware_entities.MessageTarget, 0),
		NonImmune: make([]middleware_entities.MessageTarget, 0),
	}

	for _, cachedMsg := range cache {
		if cachedMsg.IsImmune {
			lists.Immune = append(lists.Immune, cachedMsg.Message.Target)
		} else {
			lists.NonImmune = append(lists.NonImmune, cachedMsg.Message.Target)
		}
	}

	return lists
}
