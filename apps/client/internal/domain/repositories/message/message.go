package message

import (
	"eaglechat/apps/client/internal/domain/entities"
	"time"
)

// ChatOverview is a summary of a chat for display in a list.
type ChatOverview struct {
	Partner     entities.User
	LastMessage string
	Timestamp   time.Time
	UnreadCount int
}

type MessageRepository interface {
	Save(entities.Message) error
	GetChat(entities.UserID) ([]entities.Message, error)
	GetAllChatOverviews() ([]ChatOverview, error)
	GetLastMessageIndexFromSelf(entities.UserID) (int, error)
	IncrementUnreadCount(entities.UserID) error
	ResetUnreadCount(entities.UserID) error
}
