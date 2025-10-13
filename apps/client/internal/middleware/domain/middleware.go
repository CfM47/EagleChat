package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	message_cache "eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
	user_cache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	middleware_services "eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/common/simplecrypto/rsa"
	"time"
)

type Middleware struct {
	ownPort uint16
	ownUser entities.User

	p2pConnections map[entities.UserID]middleware_entities.P2PConnection
	messageCache   message_cache.MessageCache

	p2pConnPool middleware_services.P2PConnPool

	iDManagerPool middleware_services.IDManagerPool
	knownUsers    user_cache.UserCacheRepository

	receivedMessages chan<- entities.Message

	sk rsa.PrivateKey

	messageSenderTicker *time.Ticker
	announcementTicker  *time.Ticker
	quit                chan struct{}
}

var _ services.Middleware = (*Middleware)(nil)

func (m *Middleware) Shutdown() {
	m.messageSenderTicker.Stop()
	close(m.quit)
}

func (m *Middleware) Done() <-chan struct{} {
	return m.quit
}
