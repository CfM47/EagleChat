package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eaglechat/apps/id_manager/internal/diconfig"
	"eaglechat/apps/id_manager/internal/infrastructure/http/router"

	"github.com/gin-gonic/gin"
)

func main() {
	container, err := diconfig.NewContainer()
	if err != nil {
		panic(fmt.Sprintf("error building dependencies: %v", err))
	}
	defer container.MulticastNet.Close()

	routes := []router.Route{
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
			Path:    "/pending-messages",
			Handler: container.QueryPendingMessagesHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "pending-messages",
			Handler: container.AddPendingMessagesHandler,
		},
	}

	r := gin.Default()
	router.RegisterRoutes(r, routes)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it doesn't block the shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need adding it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the requests it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
