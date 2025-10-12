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
	"eaglechat/common/multicast/implementation"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	idManagerPort = "8081"
)

func main() {
	logFile, err := setupLogger()
	if err != nil {
		log.Fatalf("failed to set up logger: %v", err)
	}
	defer logFile.Close()

	tui := tui.New()
	clientRepo, err := sqliterepository.NewSQLiteRepository("./data/client.db")
	if err != nil {
		log.Fatalf("error creating client repository: %v", err)
	}

	connector, registerer, err := buildMiddlewareDeps()
	if err != nil {
		log.Fatalf("error creating middleware dependencies: %v", err)
	}

	controller := controller.New(tui, registerer, connector, clientRepo)

	err = controller.Run()
	if err != nil {
		log.Fatalf("error running controller: %v", err)
	}
}

func setupLogger() (*os.File, error) {
	logDir := "./data/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, "client.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	return file, nil
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

func buildRegisterer() services.Registerer {
	return idmanagerregisterer.NewRegisterer(implementation.DefaultUDPAddress, idManagerPort, time.Second*10)
}
