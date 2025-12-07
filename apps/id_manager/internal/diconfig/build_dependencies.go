package diconfig

import (
	"eaglechat/apps/id_manager/internal/application/services/gossip"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"      // New import
	"eaglechat/common/simplecrypto/x509util" // New import
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

	// --- Load Cryptographic Materials ---
	privKeyPath := os.Getenv("PRIV_KEY_PATH")
	if privKeyPath == "" {
		privKeyPath = "env/private_key.pem" // Assuming default path
	}
	myPrivKeyBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key from %s: %w", privKeyPath, err)
	}
	myPrivKey, err := rsa.PrivateKeyFromBytes(myPrivKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	certPath := os.Getenv("CERT_PATH")
	if certPath == "" {
		certPath = "env/id_manager.crt" // Assuming default path for service's own cert
	}
	myCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service certificate from %s: %w", certPath, err)
	}

	caCertPath := os.Getenv("CA_CERT_PATH")
	if caCertPath == "" {
		caCertPath = "env/ca.crt" // Assuming default path for CA root cert
	}
	certVerifier, err := x509util.NewVerifier(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate verifier: %w", err)
	}
	// --- End Cryptographic Materials Loading ---

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
	gossipService := gossip.NewGossipService(
		syncDataUC,
		peerProvider,
		ownAddress,
		gossipPort,
		gossipInterval,
		myPrivKey,       // New param
		myCert,          // New param
		certVerifier,    // New param
	)

	// Initialize other use cases that need the notifier
	queryUserDataUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)
	queryPendingMessagesUC := usecases.NewQueryPendingMessagesUseCase(pendingMessagesRepo, userRepo)
	addPendingMessagesUC := usecases.NewAddPendingMessagesUseCase(pendingMessagesRepo, userRepo, gossipService)
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo, gossipService)

	// Initialize handlers
	syncDataHandler := handlers.NewSyncDataHandler(
		syncDataUC,
		myPrivKey,    // New param
		certVerifier, // New param
	)
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
