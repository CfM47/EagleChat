package main

import (
	"fmt"
	"net/http"

	"eaglechat/apps/id_manager/internal/diconfig"
	"eaglechat/apps/id_manager/internal/infrastructure/http/router"

	"github.com/gin-gonic/gin"
)

func main() {
	container, err := diconfig.NewContainer()
	if err != nil {
		panic(fmt.Sprintf("error building dependencies: %v", err))
	}

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

	// Run server
	r.Run() // 0.0.0.0:8080
}
