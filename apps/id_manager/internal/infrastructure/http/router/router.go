package router

import (
	"net/http"

	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Method  string
	Path    string
	Handler handlers.Handler
}

func RegisterRoutes(r *gin.Engine, routes []Route) {
	// Health check endpoint
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	for _, route := range routes {
		r.Handle(route.Method, route.Path, route.Handler.Handle)
	}
}
