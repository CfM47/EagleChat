package iftest

import (
	"testing"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runProfileTests(t *testing.T, factory RepoFactory) {
	t.Run("saves and retrieves an own_profile", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		privKey := getTestKeys(0)
		user := entities.NewUser("user-1", "test-user", *privKey.PublicKey(), time.Now())
		expectedProfile := entities.NewOwnProfile(user, *privKey)

		// Act
		err := repo.SaveOwnProfile(expectedProfile)
		require.NoError(t, err)

		// Assert
		actualProfile, err := repo.GetOwnProfile()
		require.NoError(t, err)

		// Compare fields individually because rsa.PublicKey and rsa.PrivateKey are complex
		assert.Equal(t, expectedProfile.User.ID, actualProfile.User.ID)
		assert.Equal(t, expectedProfile.User.Name, actualProfile.User.Name)
		assert.Equal(t, expectedProfile.PrivateKey.ToBytes(), actualProfile.PrivateKey.ToBytes())

		expectedPub, _ := expectedProfile.User.PublicKey.ToBytes()
		actualPub, _ := actualProfile.User.PublicKey.ToBytes()
		assert.Equal(t, expectedPub, actualPub)
	})

	t.Run("returns ErrOwnProfileNotSet when none is saved", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		// Act
		_, err := repo.GetOwnProfile()

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrOwnProfileNotSet)
	})

	t.Run("updates an existing profile", func(t *testing.T) {
		t.Parallel()
		// Arrange
		repo, cleanup := factory(t)
		defer cleanup()

		// Create and save initial profile
		privKey := getTestKeys(1)
		user1 := entities.NewUser("user-1", "user-one", *privKey.PublicKey(), time.Now())
		initialProfile := entities.NewOwnProfile(user1, *privKey)
		require.NoError(t, repo.SaveOwnProfile(initialProfile))

		// Create updated profile with the same ID but different name
		user2 := user1
		user2.Name = "user-one-updated"
		updatedProfile := entities.NewOwnProfile(user2, *privKey)

		// Act
		err := repo.SaveOwnProfile(updatedProfile)
		require.NoError(t, err)

		// Assert
		actualProfile, err := repo.GetOwnProfile()
		require.NoError(t, err)
		assert.Equal(t, "user-one-updated", actualProfile.User.Name)
		assert.Equal(t, initialProfile.User.ID, actualProfile.User.ID) // ID should not change
	})
}
