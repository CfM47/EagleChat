package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

func (c *Connector) Connect(listenPort uint16, user entities.User, sk rsa.PrivateKey) (Middleware, <-chan entities.Message, error) {
	iDManagerPool, err := c.iDManagerPoolBuilder(sk, c.idManagerConnectionBuilder, user.ID)
	if err != nil {
		return Middleware{}, nil, err
	}

	p2pConnPool, err := c.p2pPoolBuilder(c.p2pDialer, c.p2pListenerStarter, listenPort)
	if err != nil {
		return Middleware{}, nil, err
	}

	var messageChannel chan entities.Message

	return Middleware{
		ownPort: listenPort,

		p2pConnections: make(map[entities.UserID]middleware_entities.P2PConnection),
		messageCache:   c.messageCache,

		p2pConnPool: p2pConnPool,

		iDManagerPool: iDManagerPool,
		knownUsers:    c.knownUsers,

		receivedMessages: (chan<- entities.Message)(messageChannel),

		quit: make(chan struct{}),
	}, (<-chan entities.Message)(messageChannel), nil
}
