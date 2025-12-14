package iftest

import (
	"testing"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runUserTests(t *testing.T, factory RepoFactory) {
	t.Run("saves and retrieves a user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKey := getTestKeys(2)
		expectedUser := entities.NewUser("user-3", "Charlie", *privKey.PublicKey(), time.Now())

		// Act
		err := repo.SaveUser(expectedUser)
		require.NoError(t, err)

		// Assert
		actualUser, err := repo.GetUser(expectedUser.ID)
		require.NoError(t, err)
		assert.Equal(t, expectedUser.Name, actualUser.Name)
		assert.Equal(t, expectedUser.ID, actualUser.ID)

		expectedPub, _ := expectedUser.PublicKey.ToBytes()
		actualPub, _ := actualUser.PublicKey.ToBytes()
		assert.Equal(t, expectedPub, actualPub)
	})

	t.Run("returns ErrUserNotFound for a non-existent user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		// Act
		_, err := repo.GetUser("non-existent-id")

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrUserNotFound)
	})

	t.Run("updates an existing user's data", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKey := getTestKeys(3)
		user := entities.NewUser("user-4", "David", *privKey.PublicKey(), time.Now())
		require.NoError(t, repo.SaveUser(user))

		user.Name = "Dave"

		// Act
		err := repo.SaveUser(user)
		require.NoError(t, err)

		// Assert
		actualUser, err := repo.GetUser(user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Dave", actualUser.Name)
	})

	t.Run("implicitly saves a new user when a message is received from them", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		// userA is our own profile, userB is the new contact
		privKeyA := getTestKeys(0)
		userA := entities.NewUser("user-A", "Alice", *privKeyA.PublicKey(), time.Now())
		profileA := entities.NewOwnProfile(userA, *privKeyA)
		require.NoError(t, repo.SaveOwnProfile(profileA))

		privKeyB := getTestKeys(1)
		userB := entities.NewUser("user-B", "Bob", *privKeyB.PublicKey(), time.Now())

		// Verify userB does not exist yet
		_, err := repo.GetUser(userB.ID)
		require.ErrorIs(t, err, repositories.ErrUserNotFound)

		// Act: Save a message from the unknown userB to ourselves (userA)
		msg := entities.NewMessage(userB, userA, "Hello from a stranger!", time.Now())
		err = repo.SaveMessage(msg)
		require.NoError(t, err)

		// Assert: The repository should have automatically saved userB
		actualUserB, err := repo.GetUser(userB.ID)
		require.NoError(t, err)
		assert.Equal(t, userB.Name, actualUserB.Name)
	})
}
