package jsonmessagecache

import (
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// DeleteImmune deletes any messages that appear in toDelete.
func (c *jsonMessageCache) DeleteImmune(toDelete []middleware_entities.MessageTarget) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, err := c.readCache()
	if err != nil {
		return err
	}

	deleteSet := make(map[middleware_entities.MessageTarget]struct{})
	for _, target := range toDelete {
		deleteSet[target] = struct{}{}
	}

	for msgID := range cache {
		if _, shouldDelete := deleteSet[msgID]; shouldDelete {
			delete(cache, msgID)
		}
	}

	return c.writeCache(cache)
}
