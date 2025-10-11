package messagecache

import (
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"time"
)

const DefaultImmunityPeriod = 30 * time.Second

type MessageCache interface {
	// StoreImmune will store a message in the cache that will not be deleted until DeleteImmune
	// is called with its id
	StoreImmune(entities.PendingMessage) error

	// StoreExpiring will store a message that will be immune to deletion for the
	// specified `immunityPeriod`. If a message with the same ID already exists
	// and is permanently immune, this operation will do nothing. If the message
	// exists and is already temporarily immune, its immunity timer will be reset
	// to the new `immunityPeriod`.
	StoreExpiring(message entities.PendingMessage, immunityPeriod time.Duration) error

	// DeleteNonPending deletes non immune messages that do not appear in pendingMessageTargets
	DeleteNonPending(pendingMessageTargets []entities.MessageTarget) error

	// DeleteImmune deletes immune messages that appear in toDelete
	DeleteImmune(toDelete []entities.MessageTarget) error

	GetAll() []entities.PendingMessage

	GetByTargetId(entities.MessageTarget) []entities.PendingMessage
}
