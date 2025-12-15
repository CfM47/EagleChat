package clientconnpool

import (
	"context"
	"sync"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
)

const (
	shutdownTimeout = 5 * time.Second
	requestTimeout  = 15 * time.Second
)

type clientConnPoolImpl struct {
	ownProfile entities.OwnProfile
	listenPort uint16
	messages   chan []middleware_entities.PendingMessage

	closeOnce sync.Once
	quit      chan struct{}
	done      chan struct{}
}

var _ services.ClientConnPool = (*clientConnPoolImpl)(nil)

// Message implements services.ClientConnPool.
func (c *clientConnPoolImpl) Message(ctx context.Context, messages []middleware_entities.PendingMessage, target middleware_entities.UserData) error {
	if target.IP == nil {
		return services.ErrInvalidIP
	}

	return c.sendMessageRequest(ctx, messages, *target.IP, target.PublicKey)
}

// Receive implements services.ClientConnPool.
func (c *clientConnPoolImpl) Receive() <-chan []middleware_entities.PendingMessage {
	return c.messages
}

// Done implements services.ClientConnPool.
func (c *clientConnPoolImpl) Done() <-chan struct{} {
	return c.done
}

// Close implements services.ClientConnPool.
func (c *clientConnPoolImpl) Close() error {
	c.closeOnce.Do(func() {
		close(c.quit)
	})
	<-c.done
	return nil
}
