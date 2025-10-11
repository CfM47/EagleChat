package messagecache

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"time"
)

const DefaultImmunityPeriod = 30 * time.Second

type PendingMessageTargetLists struct {
	Immune    []middleware_entities.MessageTarget
	NonImmune []middleware_entities.MessageTarget
}

type MessageCache interface {
	// StoreImmune will store a message in the cache that will not be deleted until DeleteImmune
	// is called with its id
	StoreImmune(middleware_entities.PendingMessage) error

	// StoreExpiring will store a message that will be immune to deletion for the
	// specified `immunityPeriod`. If a message with the same ID already exists
	// and is permanently immune, this operation will do nothing. If the message
	// exists and is already temporarily immune, its immunity timer will be reset
	// to the new `immunityPeriod`.
	StoreExpiring(message middleware_entities.PendingMessage, immunityPeriod time.Duration) error

	// DeleteNonPending deletes non immune messages that do not appear in pendingMessageTargets
	DeleteNonPending(pendingMessageTargets []middleware_entities.MessageTarget) error

	// DeleteImmune deletes any messages that appear in toDelete
	DeleteImmune(toDelete []middleware_entities.MessageTarget) error

	GetAll() []middleware_entities.PendingMessage

	GetTargets() PendingMessageTargetLists

	GetByTargetId(entities.UserID) []middleware_entities.PendingMessage
}
