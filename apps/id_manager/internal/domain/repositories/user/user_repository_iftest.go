package user

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"fmt"
	"net"
	"testing"
	"time"

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

	t.Run("updates LastSeen when IP is updated", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		require.NoError(t, repo.Save(user))

		// Get initial LastSeen
		initialUser, err := repo.FindByID(user.ID)
		require.NoError(t, err)
		initialLastSeen := initialUser.LastSeen

		// Act
		newIP := net.ParseIP("192.0.2.1")
		err = repo.UpdateIP(user.ID, newIP)
		require.NoError(t, err)

		// Assert
		updatedUser, err := repo.FindByID(user.ID)
		require.NoError(t, err)
		assert.True(t, updatedUser.LastSeen.After(initialLastSeen) || updatedUser.LastSeen.Equal(initialLastSeen))
		assert.True(t, newIP.Equal(*updatedUser.IP))
	})

	t.Run("returns nil IP when IP is expired", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		require.NoError(t, repo.Save(user))

		// Set an IP with an expired LastSeen
		ip := net.ParseIP("192.0.2.1")
		user.IP = &ip
		user.LastSeen = entities.NewUser("", "", "").LastSeen.Add(-entities.IPExpirationDuration - 1*time.Hour)
		require.NoError(t, repo.Save(user))

		// Act
		foundUser, err := repo.FindByID(user.ID)

		// Assert
		require.NoError(t, err)
		assert.Nil(t, foundUser.IP, "IP should be nil when expired")
	})

	t.Run("keeps IP when not expired", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		user := mockUser()
		require.NoError(t, repo.Save(user))

		// Set an IP with a recent LastSeen
		ip := net.ParseIP("192.0.2.1")
		err := repo.UpdateIP(user.ID, ip)
		require.NoError(t, err)

		// Act
		foundUser, err := repo.FindByID(user.ID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, foundUser.IP, "IP should not be nil when not expired")
		assert.True(t, ip.Equal(*foundUser.IP))
	})

	t.Run("returns nil IP for expired IPs in FindAll", func(t *testing.T) {
		// Arrange
		repo, cleanup := repoFactory(t)
		defer cleanup()
		
		// Create user with expired IP
		user1 := mockUser()
		ip1 := net.ParseIP("192.0.2.1")
		user1.IP = &ip1
		user1.LastSeen = entities.NewUser("", "", "").LastSeen.Add(-entities.IPExpirationDuration - 1*time.Hour)
		require.NoError(t, repo.Save(user1))

		// Create user with valid IP
		user2 := mockUser()
		require.NoError(t, repo.Save(user2))
		ip2 := net.ParseIP("192.0.2.2")
		require.NoError(t, repo.UpdateIP(user2.ID, ip2))

		// Act
		users, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		assert.Len(t, users, 2)

		// Find each user and check their IPs
		var foundUser1, foundUser2 *entities.User
		for _, u := range users {
			if u.ID == user1.ID {
				foundUser1 = u
			} else if u.ID == user2.ID {
				foundUser2 = u
			}
		}

		require.NotNil(t, foundUser1)
		require.NotNil(t, foundUser2)
		assert.Nil(t, foundUser1.IP, "User1 IP should be nil (expired)")
		assert.NotNil(t, foundUser2.IP, "User2 IP should not be nil (not expired)")
		assert.True(t, ip2.Equal(*foundUser2.IP))
	})
}
