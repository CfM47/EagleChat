package jsonmessagecache

import (
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// DeleteNonPending deletes non-immune messages that do not appear in pendingMessageTargets
// and whose grace period has expired.
func (c *jsonMessageCache) DeleteNonPending(pendingMessageTargets []middleware_entities.MessageTarget) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, err := c.readCache()
	if err != nil {
		return err
	}

	pendingSet := make(map[middleware_entities.MessageTarget]struct{})
	for _, target := range pendingMessageTargets {
		pendingSet[target] = struct{}{}
	}

	now := time.Now()
	for msgID, cachedMsg := range cache {
		// Condition 1: Message must not be permanently immune.
		if cachedMsg.IsImmune {
			continue
		}

		// Condition 2: Grace period must be over.
		if now.Before(cachedMsg.GracePeriodEnd) {
			continue
		}

		// Condition 3: Message must not be in the set of currently pending targets.
		if _, isPending := pendingSet[msgID]; isPending {
			continue
		}

		// If all conditions are met, delete the message.
		delete(cache, msgID)
	}

	return c.writeCache(cache)
}
