package sqliterepository

import (
	"fmt"
	"testing"

	"eaglechat/apps/client/internal/domain/repositories"
	"eaglechat/apps/client/internal/domain/repositories/iftest"
	"github.com/stretchr/testify/require"
)

func TestSQLiteRepository(t *testing.T) {
	factory := func(t *testing.T) (repositories.ClientRepository, func()) {
		// Arrange: Create a temporary in-memory DB for each test to ensure isolation.
		// The `cache=shared` is important to allow multiple connections from the pool
		// to access the same in-memory database.
		dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
		repo, err := NewSQLiteRepository(dbPath)
		require.NoError(t, err)

		// The sql.DB object handles its own cleanup, so the teardown function is empty.
		cleanup := func() {}

		return repo, cleanup
	}

	// Act: Run the entire interface test suite against the SQLite implementation.
	iftest.RunClientRepositoryTests(t, factory)
}
