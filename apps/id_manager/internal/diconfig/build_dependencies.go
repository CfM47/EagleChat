package diconfig

import (
	"eaglechat/apps/id_manager/internal/application/services/gossip"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"
	"eaglechat/common/simplecrypto/x509util"
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

	// --- Load CA Cryptographic Materials ---
	caPrivKeyPath := os.Getenv("CA_KEY_PATH")
	if caPrivKeyPath == "" {
		caPrivKeyPath = "ca.key"
	}
	caKeyPEM, err := os.ReadFile(caPrivKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA private key from %s: %w", caPrivKeyPath, err)
	}

	caCertPath := os.Getenv("CA_CERT_PATH")
	if caCertPath == "" {
		caCertPath = "ca.crt"
	}
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate from %s: %w", caCertPath, err)
	}
	certVerifier, err := x509util.NewVerifier(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate verifier: %w", err)
	}
	// --- End CA Cryptographic Materials Loading ---

	// --- Generate ID Manager Key Pair and Certificate at Runtime ---
	myPrivKey, myPubKey, err := rsa.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID Manager key pair: %w", err)
	}

	myCert, err := x509util.CreateAndSignCertificate(caCertPEM, caKeyPEM, myPubKey, commonName)
	if err != nil {
		return nil, fmt.Errorf("failed to create and sign ID Manager certificate: %w", err)
	}
	// --- End Runtime Generation ---

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
		myPrivKey,
		myCert,
		certVerifier,
	)

	// Initialize other use cases that need the notifier
	queryUserDataUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)
	queryPendingMessagesUC := usecases.NewQueryPendingMessagesUseCase(pendingMessagesRepo, userRepo)
	addPendingMessagesUC := usecases.NewAddPendingMessagesUseCase(pendingMessagesRepo, userRepo, gossipService)
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo, gossipService)

	// Initialize handlers
	syncDataHandler := handlers.NewSyncDataHandler(syncDataUC, myPrivKey, certVerifier)
	notifyUpdateHandler := handlers.NewNotifyUpdateHandler(gossipService)
	getRandomUsersHandler := handlers.NewGetRandomUsersHandler(getRandomUsersUC)
	queryUserHandler := handlers.NewQueryUserDataHandler(queryUserDataUC)
	queryPendingMessagesHandler := handlers.NewQueryPendingMessagesHandler(queryPendingMessagesUC)
	addPendingMessagesHandler := handlers.NewAddPendingMessagesHandler(addPendingMessagesUC)
	registerUserHandler := handlers.NewRegisterUserHandler(registerUserUC)
	pubKeyHandler := handlers.NewPubKeyHandler(myPrivKey, myCert)

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
