package idmanagerregisterer

import (
	"context"
	"fmt"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerverifier"
	"eaglechat/common/ezlog"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"
)

// registererImpl implements the services.Registerer interface.
type registererImpl struct {
	idManagerPort       uint16
	registrationTimeout time.Duration
	CAPubkey            *rsa.PublicKey
}

// NewRegisterer creates a new Registerer.
func NewRegisterer(idManagerPort uint16, timeout time.Duration, CAPubkey *rsa.PublicKey) services.Registerer {
	return &registererImpl{
		idManagerPort:       idManagerPort,
		registrationTimeout: timeout,
		CAPubkey:            CAPubkey,
	}
}

// Register orchestrates the discovery and HTTP registration process.
func (r *registererImpl) Register(ctx context.Context, username string, sk rsa.PrivateKey, CAPubkey *rsa.PublicKey) (entities.User, error) {
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

		// FIXME: make request/response use public key verification and for encryption
		_, client, err := idmanagerverifier.VerifyIDManager(ctx, ip, r.idManagerPort, r.CAPubkey)
		if err != nil {
			ezlog.Log(ctx).Warnf("ID manager failed to provide valid certificate: %v", err)
			continue
		}

		user, err := r.requestRegistration(ctx, username, sk.PublicKey(), ip.String(), client)
		if err != nil {
			ezlog.Log(ctx).Errorf("Failed to register with ID Manager at %s: %v", ip.String(), err)
			continue
		}

		return user, nil
	}

	ezlog.Log(ctx).Error("Failed to register with any discovered ID Manager")

	return entities.User{}, fmt.Errorf("failed to register with any discovered ID Manager")
}
