package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/infrastructure/http/handlers"
	persistence "eaglechat/apps/id_manager/internal/infrastructure/persistence/json"
)

func main() {
	// Persistence route definitions
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(fmt.Sprintf("error creando directorio de datos: %v", err))
	}

	// Persistence route files
	userFile := filepath.Join(dataDir, "users.json")

	// Initialize repositories
	repo := persistence.NewJSONUserRepository(userFile)

	// Initialize use cases
	queryUC := usecases.NewQueryUserDataUseCase(repo)

	// Initialize handlers
	queryHandler := handlers.NewQueryUserDataHandler(queryUC)

	// Initialize gin server
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Initialize endpoints
	r.GET("/users", func(c *gin.Context) {
		queryHandler.Handle(c)
	})

	// Run server
	r.Run() // 0.0.0.0:8080
}
