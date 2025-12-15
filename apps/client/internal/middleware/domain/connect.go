package middleware

import (
	"context"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/services"
	message_cache "eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
	user_cache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	middleware_services "eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/common/clock"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

type Connector struct {
	messageCache message_cache.MessageCache

	clientConnPoolBuilder middleware_services.ClientConnPoolBuilder

	iDManagerPoolBuilder middleware_services.IDManagerPoolBuilder

	knownUsers user_cache.UserCacheRepository
}

var _ services.Connector = (*Connector)(nil)

func NewConnector(
	messageCache message_cache.MessageCache,

	clientConnPoolBuilder middleware_services.ClientConnPoolBuilder,

	idManagerPoolBuilder middleware_services.IDManagerPoolBuilder,

	knownUsers user_cache.UserCacheRepository,
) Connector {
	return Connector{
		messageCache: messageCache,

		clientConnPoolBuilder: clientConnPoolBuilder,

		iDManagerPoolBuilder: idManagerPoolBuilder,

		knownUsers: knownUsers,
	}
}

func (c Connector) Connect(ctx context.Context, listenPort uint16, ownProfile entities.OwnProfile, CAPubkey *rsa.PublicKey) (services.Middleware, <-chan entities.Message, error) {
	ezlog.Log(ctx).Infof("Connecting as user %s on port %d", ownProfile.User.Name, listenPort)
	defer ezlog.Log(ctx).Info("Connector finished")

	ezlog.Log(ctx).Info("Building ID Manager Pool")
	iDManagerPool, err := c.iDManagerPoolBuilder.Build(ctx, ownProfile, CAPubkey)
	if err != nil {
		return &Middleware{}, nil, err
	}

	ezlog.Log(ctx).Info("Building P2P Connection Pool")
	clientConnPool, err := c.clientConnPoolBuilder.Build(ctx, ownProfile, listenPort)
	if err != nil {
		return &Middleware{}, nil, err
	}

	messageChannel := make(chan entities.Message)

	m := Middleware{
		ownPort:    listenPort,
		ownProfile: ownProfile,

		messageCache: c.messageCache,

		clientConnPool: clientConnPool,
		iDManagerPool:  iDManagerPool,

		knownUsers: c.knownUsers,

		receivedMessages: (chan<- entities.Message)(messageChannel),

		clock: clock.NewClock(),

		quit: make(chan struct{}),
	}

	receiverCtx := ezlog.NewLoggerContext("message-receiver")
	go m.messageReceiver(receiverCtx)

	senderCtx := ezlog.NewLoggerContext("message-sender")
	go m.messageSender(senderCtx)

	announcerCtx := ezlog.NewLoggerContext("presence-announcer")
	go m.presenceAnnouncer(announcerCtx)

	synchronizerCtx := ezlog.NewLoggerContext("time-synchronizer")
	go m.timeSynchronizer(synchronizerCtx)

	return &m, (<-chan entities.Message)(messageChannel), nil
}
