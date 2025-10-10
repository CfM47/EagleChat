package messagecache

import (
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"time"
)

type MessageCache interface {
	// StoreImmune will store a message in the cache that will not be deleted until DeleteImmune
	// is called with its id
	StoreImmune(entities.PendingMessage) error

	// StoreExpiring will store a message in the cache that will be immune to deletion from DeleteNonPending
	// for `immunityPeriod`, and will then available for deletion by DeleteNonPending
	StoreExpiring(message entities.PendingMessage, immunityPeriod time.Duration) error

	// DeleteNonPending deletes non immune messages that do not appear in pendingMessageTargets
	DeleteNonPending(pendingMessageTargets []entities.MessageTarget) error

	// DeleteImmune deletes immune messages that appear in toDelete
	DeleteImmune(toDelete []entities.MessageTarget) error

	GetAll() []entities.PendingMessage

	GetByTargetId(entities.MessageTarget) []entities.PendingMessage
}
