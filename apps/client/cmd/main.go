package main

import (
	"eaglechat/apps/client/internal/domain/services"
	sqliterepository "eaglechat/apps/client/internal/infrastructure/repositories/sqlite"
	middleware "eaglechat/apps/client/internal/middleware/domain"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerconn"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerregisterer"
	jsonmessagecache "eaglechat/apps/client/internal/middleware/infrastructure/messagecache/json"
	"eaglechat/apps/client/internal/middleware/infrastructure/p2pconn"
	"eaglechat/apps/client/internal/middleware/infrastructure/p2pconnpool"
	jsonusercache "eaglechat/apps/client/internal/middleware/infrastructure/usercache/json"
	"eaglechat/apps/client/internal/tui"
	"eaglechat/apps/client/internal/ui/controller"
	"eaglechat/common/ezlog"
	"eaglechat/common/ezlog/ezzerolog"
	"eaglechat/common/multicast/implementation"
	"log"
	"os"
	"time"
)

const (
	idManagerPort = "8080"
)

func main() {
	setupFallbackLogger()

	log.Print("Setting logger factory...")
	ezlog.SetLoggerFactory(ezzerolog.NewFactory())

	ctx := ezlog.NewLoggerContext("main")

	ezlog.Log(ctx).Info("Logger initialized")

	tui := tui.New()
	clientRepo, err := sqliterepository.NewSQLiteRepository("./data/client.db")
	if err != nil {
		ezlog.Log(ctx).Errorf("error creating client repository: %v", err)
		panic(err)
	}

	connector, registerer, err := buildMiddlewareDeps()
	if err != nil {
		ezlog.Log(ctx).Errorf("error creating middleware dependencies: %v", err)
		panic(err)
	}

	controller := controller.New(tui, registerer, connector, clientRepo)

	err = controller.Run()
	if err != nil {
		ezlog.Log(ctx).Errorf("error running controller: %v", err)
		panic(err)
	}
}

func buildMiddlewareDeps() (services.Connector, services.Registerer, error) {
	connector, err := buildConnector()
	if err != nil {
		return nil, nil, err
	}

	registerer := buildRegisterer()

	return connector, registerer, nil
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

func setupFallbackLogger() {
	os.MkdirAll("/data", 0755)
	file, err := os.OpenFile("/data/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}

func buildRegisterer() services.Registerer {
	return idmanagerregisterer.NewRegisterer(implementation.DefaultUDPAddress, idManagerPort, time.Second*10)
}
