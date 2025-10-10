package user

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockUser() *entities.User {
	id := entities.NewUUID()
	return entities.NewUser(id, fmt.Sprintf("user-%s", id), fmt.Sprintf("pk-%s", id))
}

func RunUserRepositoryTests(t *testing.T, repoFactory func(t *testing.T) (UserRepository, func())) {
	t.Helper()

	t.Run("saves a user", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()

		// Act
		err := repo.Save(user)

		// Assert
		require.NoError(t, err)
	})

	t.Run("finds a user by ID", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		repo.Save(user) // prerequisite

		// Act
		foundUser, err := repo.FindByID(user.ID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, user, foundUser)
	})

	t.Run("finds all users", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		repo.Save(user) // prerequisite

		// Act
		users, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, user, users[0])
	})

	t.Run("updates a user's IP", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		repo.Save(user) // prerequisite
		ip := net.ParseIP("192.0.2.1")

		// Act
		err := repo.UpdateIP(user.ID, ip)

		// Assert
		require.NoError(t, err)
		updatedUser, err := repo.FindByID(user.ID)
		require.NoError(t, err)
		assert.True(t, ip.Equal(*updatedUser.IP))
	})

	t.Run("deletes a user", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		repo.Save(user) // prerequisite

		// Act
		err := repo.Delete(user.ID)

		// Assert
		require.NoError(t, err)
		_, err = repo.FindByID(user.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("returns no users after deletion", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		repo.Save(user)      // prerequisite
		repo.Delete(user.ID) // prerequisite

		// Act
		users, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		assert.Len(t, users, 0)
	})

	t.Run("finds all users with multiple users", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user1 := mockUser()
		user2 := mockUser()
		repo.Save(user1)
		repo.Save(user2)

		// Act
		users, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.ElementsMatch(t, []*entities.User{user1, user2}, users)
	})

	t.Run("deletes only the target user", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user1 := mockUser()
		user2 := mockUser()
		user3 := mockUser()
		require.NoError(t, repo.Save(user1))
		require.NoError(t, repo.Save(user2))
		require.NoError(t, repo.Save(user3))

		// Act
		err := repo.Delete(user2.ID)

		// Assert
		require.NoError(t, err)

		users, err := repo.FindAll()
		require.NoError(t, err)
		assert.ElementsMatch(t, []*entities.User{user1, user3}, users)

		_, err = repo.FindByID(user2.ID)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("does not modify user data on insertion", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user1 := mockUser()
		user2 := mockUser()
		originalUsers := []*entities.User{user1, user2}

		// Act
		require.NoError(t, repo.Save(user1))
		require.NoError(t, repo.Save(user2))

		// Assert
		users, err := repo.FindAll()
		require.NoError(t, err)
		assert.ElementsMatch(t, originalUsers, users)
	})

	t.Run("does not modify other users on update", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user1 := mockUser()
		user2 := mockUser()
		user3 := mockUser()
		require.NoError(t, repo.Save(user1))
		require.NoError(t, repo.Save(user2))
		require.NoError(t, repo.Save(user3))

		user1BeforeUpdate, err := repo.FindByID(user1.ID)
		require.NoError(t, err)
		user3BeforeUpdate, err := repo.FindByID(user3.ID)
		require.NoError(t, err)

		// Act
		newIP := net.ParseIP("1.2.3.4")
		err = repo.UpdateIP(user2.ID, newIP)
		require.NoError(t, err)

		// Assert
		user1AfterUpdate, err := repo.FindByID(user1.ID)
		require.NoError(t, err)
		user2AfterUpdate, err := repo.FindByID(user2.ID)
		require.NoError(t, err)
		user3AfterUpdate, err := repo.FindByID(user3.ID)
		require.NoError(t, err)

		assert.Equal(t, user1BeforeUpdate, user1AfterUpdate)
		assert.Equal(t, user3BeforeUpdate, user3AfterUpdate)
		assert.True(t, newIP.Equal(*user2AfterUpdate.IP))
	})
}
