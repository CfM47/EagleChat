package repositories

import (
	"eaglechat/apps/client/internal/domain/entities"
	"errors"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrOwnProfileNotSet = errors.New("own profile has not been set")
)

// ChatOverview is a summary of a chat for display in a list.
type ChatOverview struct {
	Partner     entities.User
	LastMessage *string
	Timestamp   time.Time
	UnreadCount int
}

// ClientRepository manages the storage of chats and users for the client
//
// It supports an unread message count for each chat, with the ability
// to increment it by one, or reset it to 0. It is expected to be thread-safe
type ClientRepository interface {
	// Save stores a message in the repository, except if the message has the same
	// (ID, Sender.ID) pair as some other message in the repository, if this happens,
	// the message is assumed to be a duplicate and ignored
	SaveMessage(entities.Message) error

	// GetChat returns all messages for a given chat, sorted by creation time
	GetChat(entities.UserID) ([]entities.Message, error)

	// GetAllChatOverviews returns the corresponding ChatOverview object for each chat
	// in the repository, including empty chats (with UnreadCount = 0 and LastMessage = nil).
	GetAllChatOverviews() ([]ChatOverview, error)

	IncrementUnreadCount(entities.UserID) error
	ResetUnreadCount(entities.UserID) error

	// Save saves public data about a user.
	SaveUser(user entities.User) error
	// Get retrieves public data about a user.
	GetUser(userID entities.UserID) (entities.User, error)

	// SaveOwnProfile saves the full user profile, including the private key, to local storage.
	SaveOwnProfile(profile entities.OwnProfile) error
	// GetOwnProfile retrieves the full user profile from local storage.
	GetOwnProfile() (entities.OwnProfile, error)
}
