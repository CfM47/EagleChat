package idmanagerregisterer

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/common/multicast/implementation"
	"eaglechat/common/simplecrypto/rsa"
	"fmt"
	"log"
	"time"

	multicast "eaglechat/common/multicast/interface"
)

// registererImpl implements the services.Registerer interface.
type registererImpl struct {
	multicastAddress    string
	idManagerPort       string
	registrationTimeout time.Duration
}

// NewRegisterer creates a new Registerer.
func NewRegisterer(multicastAddress, idManagerPort string, timeout time.Duration) services.Registerer {
	return &registererImpl{
		multicastAddress:    multicastAddress,
		idManagerPort:       idManagerPort,
		registrationTimeout: timeout,
	}
}

// Register orchestrates the discovery and HTTP registration process.
func (r *registererImpl) Register(username string, sk rsa.PrivateKey) (entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.registrationTimeout)
	defer cancel()

	multicastNet, err := implementation.New(r.multicastAddress)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to initialize multicast network: %w", err)
	}
	defer multicastNet.Close()

	idManagerChan := make(chan multicast.IDManagerMessage, 1)

	// 2: maximum amount of goroutines
	errChan := make(chan error, 2)

	go r.broadcastLoop(ctx, multicastNet, errChan)
	go r.listenForIDManager(ctx, multicastNet, idManagerChan, errChan)

	select {
	case <-ctx.Done():
		return entities.User{}, fmt.Errorf("registration timed out: %w", ctx.Err())
	case err := <-errChan:
		return entities.User{}, err
	case idManager := <-idManagerChan:
		log.Printf("Discovered ID Manager %s at %s", idManager.ID, idManager.IP)
		return r.performHTTPRequest(ctx, username, sk.PublicKey(), idManager.IP)
	}
}
