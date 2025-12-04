package diconfig

import (
	"eaglechat/apps/client/internal/domain/services"
	middleware "eaglechat/apps/client/internal/middleware/domain"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerconn"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerregisterer"
	jsonmessagecache "eaglechat/apps/client/internal/middleware/infrastructure/messagecache/json"
	"eaglechat/apps/client/internal/middleware/infrastructure/p2pconn"
	"eaglechat/apps/client/internal/middleware/infrastructure/p2pconnpool"
	jsonusercache "eaglechat/apps/client/internal/middleware/infrastructure/usercache/json"
	"eaglechat/apps/client/internal/tui"
	"eaglechat/common/multicast/implementation"
	"time"
)

func BuildMiddlewareDeps(idManagerPort string) (*tui.TUI, services.Connector, services.Registerer, error) {
	// Build tui
	tui := tui.New()

	// Build connector
	connector, err := buildConnector()
	if err != nil {
		return tui, nil, nil, err
	}

	// Build registerer
	registerer := buildRegisterer(idManagerPort)

	return tui, connector, registerer, nil
}

func buildConnector() (services.Connector, error) {
	messageCache, err := jsonmessagecache.NewJSONMessageCache("./data/message_cache.json")
	if err != nil {
		return nil, err
	}

	p2pPoolBuilder := p2pconnpool.BuildP2PConnPool
	p2pDialer := p2pconn.Dial
	p2pListenerStater := p2pconn.StartListener

	idManagerConnectionBuilder := idmanagerconn.BuildIDManagerConnection
	idManagerPoolBuilder := idmanagerpool.BuildIDManagerPool

	userCache, err := jsonusercache.NewJSONUserCache("./data/user_cache.json", time.Second*10)
	if err != nil {
		return nil, err
	}

	return middleware.NewConnector(
		messageCache,

		p2pPoolBuilder,
		p2pDialer,
		p2pListenerStater,

		idManagerConnectionBuilder,
		idManagerPoolBuilder,

		userCache,
	), nil
}

func buildRegisterer(idManagerPort string) services.Registerer {
	return idmanagerregisterer.NewRegisterer(implementation.DefaultUDPAddress, idManagerPort, time.Second*10)
}
