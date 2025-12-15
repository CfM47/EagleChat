package diconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"eaglechat/apps/id_manager/internal/application/services/gossip"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
	"eaglechat/common/clock"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"
)

type Container struct {
	GetTimeHandler        handlers.Handler
	QueryUserHandler      handlers.Handler
	GetRandomUsersHandler handlers.Handler
	RegisterUserHandler   handlers.Handler
	SyncDataHandler       handlers.Handler
	NotifyUpdateHandler   handlers.Handler
	PubKeyHandler         handlers.Handler
	AnnounceHandler       handlers.Handler
	GossipService         *gossip.GossipService
}

func NewContainer() (*Container, error) {
	// Persistence route definitions
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("error creating data directory: %w", err)
	}

	// Config
	gossipInterval := 10 * time.Second
	expirationDuration := 30 * time.Second

	// --- Load CA Public Key ---
	caPubKeyPath := os.Getenv("CA_PUBLIC_KEY_PATH")
	if caPubKeyPath == "" {
		caPubKeyPath = "env/ca_public_key.pem"
	}
	caPubKey, err := rsa.PublicKeyFromFile(caPubKeyPath)
	if err != nil {
		return nil, err
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

	// Initialize clock
	clock := clock.NewClock()

	// Initialize repositories
	userRepo := persistence.NewJSONUserRepository(filepath.Join(dataDir, "users.json"), clock)

	// Initialize discovery and gossip services
	dnsDiscovery := ns.NewDNSDiscovery()
	cachedDiscovery := ns.NewCachingDiscovery(dnsDiscovery, filepath.Join(dataDir, "ip_cache.json"))
	peerProvider := gossip.NewDnsPeerProvider(cachedDiscovery)

	ownAddress, err := dnsDiscovery.GetOwnHost()
	if err != nil {
		return nil, fmt.Errorf("failed to get own host address: %w", err)
	}
	gossipPort := "8080" // we could also make this more robust

	// Usecases need to be created before gossip service if gossip service depends on them
	syncDataUC := usecases.NewSyncDataUseCase(userRepo, clock, expirationDuration)

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
	getTimeUsecase := usecases.NewGetTimeUseCase(clock)
	queryUserDataUC := usecases.NewQueryUserDataUseCase(userRepo)
	getRandomUsersUC := usecases.NewGetRandomUsersUseCase(userRepo)
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo, gossipService, clock)
	announceUseCase := usecases.NewAnnounceUseCase(userRepo, gossipService, clock, expirationDuration)

	// Initialize handlers
	getTimeHandler := handlers.NewGetTimeHandler(getTimeUsecase)
	syncDataHandler := handlers.NewSyncDataHandler(syncDataUC, myPrivKey, caPubKey)
	notifyUpdateHandler := handlers.NewNotifyUpdateHandler(gossipService)
	getRandomUsersHandler := handlers.NewGetRandomUsersHandler(getRandomUsersUC)
	queryUserHandler := handlers.NewQueryUserDataHandler(queryUserDataUC)
	registerUserHandler := handlers.NewRegisterUserHandler(registerUserUC)
	pubKeyHandler := handlers.NewPubKeyHandler(myPrivKey, mySignature)
	AnnounceHandler := handlers.NewAnnounceHandler(announceUseCase)

	return &Container{
		GetTimeHandler:        getTimeHandler,
		QueryUserHandler:      queryUserHandler,
		GetRandomUsersHandler: getRandomUsersHandler,
		RegisterUserHandler:   registerUserHandler,
		SyncDataHandler:       syncDataHandler,
		NotifyUpdateHandler:   notifyUpdateHandler,
		PubKeyHandler:         pubKeyHandler,
		AnnounceHandler:       AnnounceHandler,
		GossipService:         gossipService,
	}, nil
}
