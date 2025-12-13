package idmanagerpool

import (
	"context"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/idmanagerconn"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/repositories"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

type idManagerPoolBuilderImpl struct{}

var _ services.IDManagerPoolBuilder = (*idManagerPoolBuilderImpl)(nil)

func NewIDManagerPoolBuilder() services.IDManagerPoolBuilder {
	return &idManagerPoolBuilderImpl{}
}

// Build implements services.IDManagerPoolBuilder.
func (i *idManagerPoolBuilderImpl) Build(ctx context.Context, ownProfile entities.OwnProfile, CAPubkey rsa.PublicKey) (services.IDManagerPool, error) {
	ezlog.Log(ctx).Info("Building ID manager pool")

	repo := repositories.NewInMemoryIDManagerRepository(ExpirationTime)

	pool := &idManagerPoolImpl{
		repository: repo,
		ownProfile: ownProfile,
		connector:  idmanagerconn.NewIDManagerConnector(ownProfile, CAPubkey),
		CAPubkey:   CAPubkey,
		quitChan:   make(chan struct{}),
		doneChan:   make(chan struct{}),
	}

	pollCtx := ezlog.WithNewLogger(ctx, "id-manager-pool-poll-loop")
	go pool.pollDNSLoop(pollCtx)

	return pool, nil
}
