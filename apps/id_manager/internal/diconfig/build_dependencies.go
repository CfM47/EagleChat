package diconfig

import (
	"eaglechat/apps/id_manager/internal/application/services/gossip"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
	"eaglechat/common/ns"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Container struct {
	QueryUserHandler            handlers.Handler
	GetRandomUsersHandler       handlers.Handler
	QueryPendingMessagesHandler handlers.Handler
	AddPendingMessagesHandler   handlers.Handler
	RegisterUserHandler         handlers.Handler
	SyncDataHandler             handlers.Handler
	NotifyUpdateHandler         handlers.Handler
	GossipService               *gossip.GossipService
}

func NewContainer() (*Container, error) {
	// Persistence route definitions
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating data directory: %w", err)
	}

	// Config
	ownAddress := "localhost:8080"     // Hardcoded for now
	gossipPort := "8080"               // Hardcoded for now
	gossipInterval := 10 * time.Second // Hardcoded for now

	// Initialize repositories
	userRepo := persistence.NewJSONUserRepository(filepath.Join(dataDir, "users.json"))
	pendingMessagesRepo := persistence.NewJSONPendingMessageRepository(filepath.Join(dataDir, "pending_messages.json"))

	// Initialize discovery and gossip services
	dnsDiscovery := ns.NewDNSDiscovery()
	cachedDiscovery := ns.NewCachingDiscovery(dnsDiscovery, filepath.Join(dataDir, "ip_cache.json"))
	peerProvider := gossip.NewDnsPeerProvider(cachedDiscovery)

	// Usecases need to be created before gossip service if gossip service depends on them
	syncDataUC := usecases.NewSyncDataUseCase(userRepo, pendingMessagesRepo)

	// Create gossip service
	gossipService := gossip.NewGossipService(syncDataUC, peerProvider, ownAddress, gossipPort, gossipInterval)

	// Initialize other use cases that need the notifier
	queryUserDataUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)
	queryPendingMessagesUC := usecases.NewQueryPendingMessagesUseCase(pendingMessagesRepo, userRepo)
	addPendingMessagesUC := usecases.NewAddPendingMessagesUseCase(pendingMessagesRepo, userRepo, gossipService)
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo, gossipService)

	// Initialize handlers
	syncDataHandler := handlers.NewSyncDataHandler(syncDataUC)
	notifyUpdateHandler := handlers.NewNotifyUpdateHandler(gossipService)
	getRandomUsersHandler := handlers.NewGetRandomUsersHandler(getRandomUsersUC)
	queryUserHandler := handlers.NewQueryUserDataHandler(queryUserDataUC)
	queryPendingMessagesHandler := handlers.NewQueryPendingMessagesHandler(queryPendingMessagesUC)
	addPendingMessagesHandler := handlers.NewAddPendingMessagesHandler(addPendingMessagesUC)
	registerUserHandler := handlers.NewRegisterUserHandler(registerUserUC)

	return &Container{
		QueryUserHandler:            queryUserHandler,
		GetRandomUsersHandler:       getRandomUsersHandler,
		QueryPendingMessagesHandler: queryPendingMessagesHandler,
		AddPendingMessagesHandler:   addPendingMessagesHandler,
		RegisterUserHandler:         registerUserHandler,
		SyncDataHandler:             syncDataHandler,
		NotifyUpdateHandler:         notifyUpdateHandler,
		GossipService:               gossipService,
	}, nil
}
