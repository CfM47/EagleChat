package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"eaglechat/apps/id_manager/internal/diconfig"
	"eaglechat/apps/id_manager/internal/infrastructure/http/router"
	"eaglechat/common/ezlog"
	"eaglechat/common/ezlog/ezzerolog"

	"github.com/gin-gonic/gin"
)

func main() {
	setupFallbackLogger()

	log.Print("Setting logger factory...")
	ezlog.SetLoggerFactory(ezzerolog.NewFactory())
	ctx := ezlog.NewLoggerContext("main")
	ezlog.Log(ctx).Info("Logger initialized")

	container, err := diconfig.NewContainer()
	if err != nil {
		panic(fmt.Sprintf("error building dependencies: %v", err))
	}

	routes := []router.Route{
		{
			Method:  http.MethodGet,
			Path:    "/time",
			Handler: container.GetTimeHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: container.QueryUserHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/random",
			Handler: container.GetRandomUsersHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users/register",
			Handler: container.RegisterUserHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/pubkey",
			Handler: container.PubKeyHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "/sync",
			Handler: container.SyncDataHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "/notify-update",
			Handler: container.NotifyUpdateHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "/announce",
			Handler: container.AnnounceHandler,
		},
	}

	r := gin.Default()
	router.RegisterRoutes(r, routes)

	// Start the gossip service
	container.GossipService.Start()

	// Run server
	r.Run() // 0.0.0.0:8080
}

func setupFallbackLogger() {
	err := os.MkdirAll("/data", 0o755)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("/data/id_manager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}
