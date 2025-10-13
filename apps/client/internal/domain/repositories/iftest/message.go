package iftest

import (
	"eaglechat/apps/client/internal/domain/entities"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runMessageTests(t *testing.T, factory RepoFactory) {
	t.Run("retrieves a generated chat in chronological order", func(t *testing.T) {
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

		generatedChat := generateChat(userA, userB, 5, 5)

		// Act
		insertChat(t, repo, generatedChat)

		// Assert
		retrievedChat, err := repo.GetChat(userB.ID)
		require.NoError(t, err)

		// Sort the generated chat
		sort.Sort(sortByCreatedDate(generatedChat))

		assert.Equal(t, len(generatedChat), len(retrievedChat))
		for i := range generatedChat {
			assert.Equal(t, generatedChat[i].ID, retrievedChat[i].ID)
		}
	})

	t.Run("correctly separates messages from multiple generated chats", func(t *testing.T) {
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
		privKeyC := getTestKeys(2)
		userC := entities.NewUser("user-C", "Charlie", *privKeyC.PublicKey())

		chatB := generateChat(userA, userB, 3, 2) // 5 messages
		chatC := generateChat(userA, userC, 4, 1) // 5 messages

		// Act
		insertChat(t, repo, chatB)
		insertChat(t, repo, chatC)

		// Assert
		retrievedB, err := repo.GetChat(userB.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedB, 5)

		retrievedC, err := repo.GetChat(userC.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedC, 5)
	})

	t.Run("ignores duplicate messages", func(t *testing.T) {
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

		msg := generateChat(userA, userB, 1, 0)[0]

		// Act
		require.NoError(t, repo.SaveMessage(msg)) // First save
		require.NoError(t, repo.SaveMessage(msg)) // Second save (should be ignored)

		// Assert
		retrievedChat, err := repo.GetChat(userB.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedChat, 1)
	})

	t.Run("overview for user with no messages is empty", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		// We need a profile for the "self" user for the repository to work correctly
		// with chat directions, even if no messages are sent from self.
		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		// This is the user we will have an empty chat with.
		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey())
		require.NoError(t, repo.SaveUser(userB))

		// Act
		overviews, err := repo.GetAllChatOverviews()
		require.NoError(t, err)

		// Assert
		require.Len(t, overviews, 1, "Expected one chat overview for the saved user")

		overview := overviews[0]
		assert.Equal(t, userB.ID, overview.Partner.ID)
		assert.Nil(t, overview.LastMessage)
		assert.Zero(t, overview.UnreadCount)
	})

	t.Run("get chat on user with no messages returns empty list", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		// We need a profile for the "self" user for the repository to work correctly
		// with chat directions, even if no messages are sent from self.
		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		// This is the user we will have an empty chat with.
		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey())
		require.NoError(t, repo.SaveUser(userB))

		// Act
		messages, err := repo.GetChat(userB.ID)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, messages, "Expected no messages for user with no chat history")
	})
}
