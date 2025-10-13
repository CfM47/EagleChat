package json_test

import (
	repository "eaglechat/apps/id_manager/internal/domain/repositories/user"
	json_user_repository "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
	"path/filepath"
	"testing"
)

func TestJsonUserRepository(t *testing.T) {
	repository.RunUserRepositoryTests(t, func(t *testing.T) (repository.UserRepository, func()) {
		t.Helper()

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_db.json")

		repo := json_user_repository.NewJSONUserRepository(filePath)

		cleanup := func() {} // No-op, since t.TempDir() handles cleanup

		return repo, cleanup
	})
}
