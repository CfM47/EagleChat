package diconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
)

type Container struct {
	QueryHandler       handlers.Handler
	GetRandomUsersHandler handlers.Handler
}

func NewContainer() (*Container, error) {
	// Persistence route definitions
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating data directory: %v", err)
	}

	// Persistence route files
	userFile := filepath.Join(dataDir, "users.json")

	// Initialize repositories
	userRepo := persistence.NewJSONUserRepository(userFile)

	// Initialize use cases
	queryUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)

	// Initialize handlers
	queryHandler := handlers.NewQueryUserDataHandler(queryUC)
	getRandomUsersHandler := handlers.NewGetRandomUsersHandler(getRandomUsersUC)

	return &Container{
		QueryHandler:       queryHandler,
		GetRandomUsersHandler: getRandomUsersHandler,
	}, nil
}
