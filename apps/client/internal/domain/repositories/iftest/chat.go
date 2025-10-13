package iftest

import (
	"eaglechat/apps/client/internal/domain/entities"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runChatTests(t *testing.T, factory RepoFactory) {
	t.Run("initially creates chat with zero unread messages", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey())

		// Act
		msg := entities.NewMessage(userB, userA, "Hello")
		require.NoError(t, repo.SaveMessage(msg))

		// Assert
		overviews, err := repo.GetAllChatOverviews()
		require.NoError(t, err)
		require.Len(t, overviews, 1)
		assert.Equal(t, 0, overviews[0].UnreadCount)
	})

	t.Run("IncrementUnreadCount correctly increases the count", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey())
		insertChat(t, repo, []entities.Message{entities.NewMessage(userB, userA, "Hi")})

		// Act
		incrementUnreadCount(t, repo, userB.ID, 3)

		// Assert
		overviews, err := repo.GetAllChatOverviews()
		require.NoError(t, err)
		require.Len(t, overviews, 1)
		assert.Equal(t, 3, overviews[0].UnreadCount)
	})

	t.Run("ResetUnreadCount correctly sets the count to zero", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey())
		insertChat(t, repo, []entities.Message{entities.NewMessage(userB, userA, "Hi")})
		incrementUnreadCount(t, repo, userB.ID, 5)

		// Act
		require.NoError(t, repo.ResetUnreadCount(userB.ID))

		// Assert
		overviews, err := repo.GetAllChatOverviews()
		require.NoError(t, err)
		require.Len(t, overviews, 1)
		assert.Equal(t, 0, overviews[0].UnreadCount)
	})

	t.Run("GetAllChatOverviews returns correct last message and timestamp", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey())

		chat := generateChat(userA, userB, 2, 3)
		insertChat(t, repo, chat)

		// Find the actual last message
		sort.Sort(sortByCreatedDate(chat))
		lastMessage := chat[len(chat)-1]

		// Act
		overviews, err := repo.GetAllChatOverviews()
		require.NoError(t, err)

		// Assert
		require.Len(t, overviews, 1)
		overview := overviews[0]
		assert.Equal(t, lastMessage.Content, *overview.LastMessage)
		assert.True(t, lastMessage.CreatedTime.Equal(overview.Timestamp))
		assert.Equal(t, userB.ID, overview.Partner.ID)
	})
}
