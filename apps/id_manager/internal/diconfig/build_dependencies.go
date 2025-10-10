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
	QueryUserHandler            handlers.Handler
	GetRandomUsersHandler       handlers.Handler
	QueryPendingMessagesHandler handlers.Handler
	AddPendingMessagesHandler   handlers.Handler
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
	pendingMessagesFile := filepath.Join(dataDir, "pending_messages.json")

	// Initialize repositories
	userRepo := persistence.NewJSONUserRepository(userFile)
	pendingMessagesRepo := persistence.NewJSONPendingMessageRepository(pendingMessagesFile)

	// Initialize use cases
	queryUserDataUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)
	queryPendingMessagesUC := usecases.NewQueryPendingMessagesUseCase(pendingMessagesRepo)
	addPendingMessagesUC := usecases.NewAddPendingMessagesUseCase(pendingMessagesRepo)

	// Initialize handlers
	getRandomUsersHandler := handlers.NewGetRandomUsersHandler(getRandomUsersUC)
	queryUserHandler := handlers.NewQueryUserDataHandler(queryUserDataUC)
	queryPendingMessagesHandler := handlers.NewQueryPendingMessagesHandler(queryPendingMessagesUC)
	addPendingMessagesHandler := handlers.NewAddPendingMessagesHandler(addPendingMessagesUC)

	return &Container{
		QueryUserHandler:            queryUserHandler,
		GetRandomUsersHandler:       getRandomUsersHandler,
		QueryPendingMessagesHandler: queryPendingMessagesHandler,
		AddPendingMessagesHandler:   addPendingMessagesHandler,
	}, nil
}
