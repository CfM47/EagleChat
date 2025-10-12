package iftest

import (
	"testing"

	"eaglechat/apps/client/internal/domain/repositories"
)

// RepoFactory defines a function that can create a new ClientRepository for testing,
// along with a cleanup function to be called after the test completes.
type RepoFactory func(t *testing.T) (repo repositories.ClientRepository, cleanup func())

// RunClientRepositoryTests runs a comprehensive suite of tests against any implementation
// of the ClientRepository interface.
func RunClientRepositoryTests(t *testing.T, factory RepoFactory) {
	t.Run("Profile", func(t *testing.T) {
		runProfileTests(t, factory)
	})
	t.Run("Users", func(t *testing.T) {
		runUserTests(t, factory)
	})
	t.Run("Messages", func(t *testing.T) {
		runMessageTests(t, factory)
	})
		t.Run("Chats", func(t *testing.T) {
		runChatTests(t, factory)
	})
}
