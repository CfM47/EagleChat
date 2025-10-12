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

type Connector struct {
	messageCache message_cache.MessageCache

	p2pPoolBuilder     middleware_services.P2PConnPoolBuilder
	p2pDialer          middleware_entities.P2PDialer
	p2pListenerStarter middleware_entities.P2PListenStarter

	idManagerConnectionBuilder middleware_entities.IDManagerConnBuilder
	iDManagerPoolBuilder       middleware_services.IDManagerPoolBuilder

	knownUsers user_cache.UserCacheRepository
}

var _ services.Connector = (*Connector)(nil)

func NewConnector(
	messageCache message_cache.MessageCache,

	p2pPoolBuilder middleware_services.P2PConnPoolBuilder,
	p2pDialer middleware_entities.P2PDialer,
	p2pListenerStarter middleware_entities.P2PListenStarter,

	idManagerConnectionBuilder middleware_entities.IDManagerConnBuilder,
	idManagerPoolBuilder middleware_services.IDManagerPoolBuilder,

	knownUsers user_cache.UserCacheRepository,
) Connector {
	return Connector{
		messageCache: messageCache,

		p2pPoolBuilder:     p2pPoolBuilder,
		p2pDialer:          p2pDialer,
		p2pListenerStarter: p2pListenerStarter,

		idManagerConnectionBuilder: idManagerConnectionBuilder,
		iDManagerPoolBuilder:       idManagerPoolBuilder,

		knownUsers: knownUsers,
	}
}

func (c Connector) Connect(listenPort uint16, user entities.User, sk rsa.PrivateKey) (services.Middleware, <-chan entities.Message, error) {
	iDManagerPool, err := c.iDManagerPoolBuilder(sk, c.idManagerConnectionBuilder, user.ID)
	if err != nil {
		return &Middleware{}, nil, err
	}

	p2pConnPool, err := c.p2pPoolBuilder(c.p2pDialer, c.p2pListenerStarter, listenPort)
	if err != nil {
		return &Middleware{}, nil, err
	}

	messageChannel := make(chan entities.Message)

	m := Middleware{
		ownPort: listenPort,
		ownUser: user,
		sk:      sk,

		p2pConnections: make(map[entities.UserID]middleware_entities.P2PConnection),
		messageCache:   c.messageCache,

		p2pConnPool: p2pConnPool,

		iDManagerPool: iDManagerPool,
		knownUsers:    c.knownUsers,

		receivedMessages: (chan<- entities.Message)(messageChannel),

		quit:                make(chan struct{}),
		messageSenderTicker: time.NewTicker(messageSenderInterval),
	}

	go m.routeIncomingMessages()
	go m.messageSender()

	return &m, (<-chan entities.Message)(messageChannel), nil
}
