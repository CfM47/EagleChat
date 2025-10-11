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
			Handler: container.QueryHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/random",
			Handler: container.GetRandomUsersHandler,
		},
	}

	r := gin.Default()
	router.RegisterRoutes(r, routes)

	// Run server
	r.Run() // 0.0.0.0:8080
}
