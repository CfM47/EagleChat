package middleware

import (
	"eaglechat/apps/client/internal/domain"
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	message_cache "eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
	user_cache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

var _ domain.Middleware = (*Middleware)(nil)

type Middleware struct {
	ownPort uint16
	ownUser entities.User

	p2pConnections map[entities.UserID]middleware_entities.P2PConnection
	messageCache   message_cache.MessageCache

	p2pConnPool services.P2PConnPool

	iDManagerPool services.IDManagerPool
	knownUsers    user_cache.UserCacheRepository

	receivedMessages chan<- entities.Message

	sk rsa.PrivateKey

	quit chan struct{}
}

type Connector struct {
	messageCache message_cache.MessageCache

	p2pPoolBuilder     services.P2PConnPoolBuilder
	p2pDialer          middleware_entities.P2PDialer
	p2pListenerStarter middleware_entities.P2PListenStarter

	idManagerConnectionBuilder middleware_entities.IDManagerConnBuilder
	iDManagerPoolBuilder       services.IDManagerPoolBuilder

	knownUsers user_cache.UserCacheRepository
}
