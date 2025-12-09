package diconfig

import (
	"eaglechat/apps/id_manager/internal/application/services/gossip"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"
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
	PubKeyHandler               handlers.Handler
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
	ownAddress := "localhost:8080"
	gossipPort := "8080"
	gossipInterval := 10 * time.Second
	commonName := os.Getenv("COMMON_NAME")
	if commonName == "" {
		commonName = ownAddress
	}

	// --- Load CA Public Key ---
	caPubKeyPath := os.Getenv("CA_PUBLIC_KEY_PATH")
	if caPubKeyPath == "" {
		caPubKeyPath = "env/ca_public_key.pem"
	}
	caPubKeyBytes, err := os.ReadFile(caPubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA public key from %s: %w", caPubKeyPath, err)
	}
	caPubKey, err := rsa.PublicKeyFromBytes(caPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA public key: %w", err)
	}

	// --- Load ID Manager's Private Key and derive Public Key ---
	idManagerPrivKeyPath := os.Getenv("ID_MANAGER_PRIV_KEY_PATH")
	if idManagerPrivKeyPath == "" {
		idManagerPrivKeyPath = "env/private_key.pem"
	}
	myPrivKeyBytes, err := os.ReadFile(idManagerPrivKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ID Manager private key from %s: %w", idManagerPrivKeyPath, err)
	}
	myPrivKey, err := rsa.PrivateKeyFromBytes(myPrivKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID Manager private key: %w", err)
	}
	myPubKey := myPrivKey.PublicKey()

	// --- Load ID Manager's Signature (signed by CA) ---
	idManagerSignaturePath := os.Getenv("ID_MANAGER_SIGNATURE_PATH")
	if idManagerSignaturePath == "" {
		idManagerSignaturePath = "env/id_manager_signature.pem"
	}
	mySignature, err := os.ReadFile(idManagerSignaturePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ID Manager signature from %s: %w", idManagerSignaturePath, err)
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
		myPubKey,
		myPrivKey,
		mySignature,
		caPubKey,
	)

	// Initialize other use cases that need the notifier
	queryUserDataUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)
	queryPendingMessagesUC := usecases.NewQueryPendingMessagesUseCase(pendingMessagesRepo, userRepo)
	addPendingMessagesUC := usecases.NewAddPendingMessagesUseCase(pendingMessagesRepo, userRepo, gossipService)
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo, gossipService)

	// Initialize handlers
	syncDataHandler := handlers.NewSyncDataHandler(syncDataUC, myPrivKey, caPubKey)
	notifyUpdateHandler := handlers.NewNotifyUpdateHandler(gossipService)
	getRandomUsersHandler := handlers.NewGetRandomUsersHandler(getRandomUsersUC)
	queryUserHandler := handlers.NewQueryUserDataHandler(queryUserDataUC)
	queryPendingMessagesHandler := handlers.NewQueryPendingMessagesHandler(queryPendingMessagesUC)
	addPendingMessagesHandler := handlers.NewAddPendingMessagesHandler(addPendingMessagesUC)
	registerUserHandler := handlers.NewRegisterUserHandler(registerUserUC)
	pubKeyHandler := handlers.NewPubKeyHandler(myPrivKey, mySignature)

	return &Container{
		QueryUserHandler:            queryUserHandler,
		GetRandomUsersHandler:       getRandomUsersHandler,
		QueryPendingMessagesHandler: queryPendingMessagesHandler,
		AddPendingMessagesHandler:   addPendingMessagesHandler,
		RegisterUserHandler:         registerUserHandler,
		SyncDataHandler:             syncDataHandler,
		NotifyUpdateHandler:         notifyUpdateHandler,
		PubKeyHandler:               pubKeyHandler,
		GossipService:               gossipService,
	}, nil
}
