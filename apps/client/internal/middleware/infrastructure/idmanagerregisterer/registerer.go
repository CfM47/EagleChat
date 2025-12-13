package idmanagerregisterer

import (
	"context"
	"fmt"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/common/ezlog"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"
)

// registererImpl implements the services.Registerer interface.
type registererImpl struct {
	multicastAddress    string
	idManagerPort       string
	registrationTimeout time.Duration
}

// NewRegisterer creates a new Registerer.
// FIXME: make request/response use public key verification and for encryption
func NewRegisterer(multicastAddress, idManagerPort string, timeout time.Duration) services.Registerer {
	return &registererImpl{
		multicastAddress:    multicastAddress,
		idManagerPort:       idManagerPort,
		registrationTimeout: timeout,
	}
}

// Register orchestrates the discovery and HTTP registration process.
func (r *registererImpl) Register(ctx context.Context, username string, sk rsa.PrivateKey, CAPubkey rsa.PublicKey) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.registrationTimeout)
	ezlog.Log(ctx).Info("Beginning registration process")

	defer cancel()

	idManagerIps, err := ns.NewDNSDiscovery().DiscoverIDManagerIPs(ctx)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to discover ID Manager IPs: %v", err)
		return entities.User{}, fmt.Errorf("failed to discover ID Manager IPs: %w", err)
	}

	for _, ip := range idManagerIps {
		ezlog.Log(ctx).Infof("Trying to register with id manager at %s", ip.String())

		user, err := r.performHTTPRequest(ctx, username, sk.PublicKey(), ip.String())
		if err != nil {
			ezlog.Log(ctx).Errorf("Failed to register with ID Manager at %s: %v", ip.String(), err)
			continue
		}

		return user, nil
	}

	ezlog.Log(ctx).Error("Failed to register with any discovered ID Manager")

	return entities.User{}, fmt.Errorf("failed to register with any discovered ID Manager")
}
