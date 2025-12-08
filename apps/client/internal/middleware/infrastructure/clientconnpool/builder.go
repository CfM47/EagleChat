package clientconnpool

import (
	"context"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
)

type clientConnBuilderImpl struct{}

var _ services.ClientConnPoolBuilder = (*clientConnBuilderImpl)(nil)

func NewClientConnPoolBuilder() services.ClientConnPoolBuilder {
	return &clientConnBuilderImpl{}
}

// Build implements services.ClientConnPoolBuilder.
func (c *clientConnBuilderImpl) Build(ctx context.Context, ownProfile entities.OwnProfile, listenPort uint16) (services.ClientConnPool, error) {
	pool := &clientConnPoolImpl{
		listenPort: listenPort,
		ownProfile: ownProfile,
		messages:   make(chan []middleware_entities.PendingMessage, 1),
		quit:       make(chan struct{}),
		done:       make(chan struct{}),
	}

	go pool.serve()

	return pool, nil
}
