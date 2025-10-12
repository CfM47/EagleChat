package idmanagerregisterer

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/common/ezlog"
	"eaglechat/common/multicast/implementation"
	"eaglechat/common/simplecrypto/rsa"
	"fmt"
	"net"
	"time"
)

// registererImpl implements the services.Registerer interface.
type registererImpl struct {
	multicastAddress    string
	idManagerPort       string
	registrationTimeout time.Duration

	foundIDManager chan struct{}
}

// NewRegisterer creates a new Registerer.
func NewRegisterer(multicastAddress, idManagerPort string, timeout time.Duration) services.Registerer {
	return &registererImpl{
		multicastAddress:    multicastAddress,
		idManagerPort:       idManagerPort,
		registrationTimeout: timeout,
	}
}

type IDManagerData struct {
	IP   net.IP
	Port uint16
}

// Register orchestrates the discovery and HTTP registration process.
func (r *registererImpl) Register(ctx context.Context, username string, sk rsa.PrivateKey) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.registrationTimeout)
	ctx = ezlog.WithComponentPrefix(ctx, "Registerer")
	ezlog.Log(ctx).Info("beginning registration process")

	defer cancel()

	multicastNet, err := implementation.New(r.multicastAddress)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to initialize multicast network: %w", err)
	}
	defer multicastNet.Close()

	idManagerChan := make(chan IDManagerData, 1)

	// 2: maximum amount of goroutines
	errChan := make(chan error, 2)

	go r.broadcastLoop(ctx, multicastNet, errChan)
	go r.listenForIDManager(ctx, multicastNet, idManagerChan, errChan)
	go r.tryDefaultIDManager(ctx, idManagerChan)

	select {
	case <-ctx.Done():
		return entities.User{}, fmt.Errorf("registration timed out: %w", ctx.Err())
	case err := <-errChan:
		return entities.User{}, err
	case idManager := <-idManagerChan:
		r.foundIDManager <- struct{}{}

		ezlog.Log(ctx).Infof("Discovered ID Manager at %s", idManager.IP)
		return r.performHTTPRequest(ctx, username, sk.PublicKey(), idManager.IP.String())
	}
}
