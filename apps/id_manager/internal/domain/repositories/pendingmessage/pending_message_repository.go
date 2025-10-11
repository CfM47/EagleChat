package pendingmessage

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"errors"
)

type PendingMessageRepository interface {
	Save(item *entities.PendingMessage) error
	FindByID(message_id string, target_id string) (*entities.PendingMessage, error)
	FindAll() ([]*entities.PendingMessage, error)
	FindByTargetID(target_id string) ([]*entities.PendingMessage, error)
	Delete(message_id string, target_id string) error
	RemoveCacher(message_id string, target_id string, cacher_id string) error
	RemoveCachersMany(message_id string, target_id string, cachers_id []string) error
}

var ErrPendingMessageNotFound = errors.New("user not found")
