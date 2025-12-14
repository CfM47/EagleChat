package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	message_cache "eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
	user_cache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	middleware_services "eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/common/clock"
)

type Middleware struct {
	ownPort    uint16
	ownProfile entities.OwnProfile

	messageCache message_cache.MessageCache

	clientConnPool middleware_services.ClientConnPool
	iDManagerPool  middleware_services.IDManagerPool

	knownUsers user_cache.UserCacheRepository

	receivedMessages chan<- entities.Message

	clock clock.Clock

	quit chan struct{}
}

var _ services.Middleware = (*Middleware)(nil)

func (m *Middleware) Shutdown() {
	close(m.quit)
}

func (m *Middleware) Done() <-chan struct{} {
	return m.quit
}
