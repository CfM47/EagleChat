package main

import (
	"eaglechat/apps/client/internal/diconfig"
	sqliterepository "eaglechat/apps/client/internal/infrastructure/repositories/sqlite"
	"eaglechat/apps/client/internal/ui/controller"
	"eaglechat/common/ezlog"
	"eaglechat/common/ezlog/ezzerolog"
	"log"
	"os"
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

	clientRepo, err := sqliterepository.NewSQLiteRepository("./data/client.db")
	if err != nil {
		ezlog.Log(ctx).Errorf("error creating client repository: %v", err)
		panic(err)
	}

	tui, connector, registerer, err := diconfig.BuildMiddlewareDeps(idManagerPort)
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

func setupFallbackLogger() {
	err := os.MkdirAll("/data", 0o755)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("/data/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}
