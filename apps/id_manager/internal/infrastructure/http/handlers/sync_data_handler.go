package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases"
	"log" // Temporarily using standard log
	"net/http"

	"github.com/gin-gonic/gin"
)

// SyncDataHandler handles GET requests for /sync endpoint.
type SyncDataHandler struct {
	syncDataUC *usecases.SyncDataUseCase
}

// NewSyncDataHandler creates a new SyncDataHandler.
func NewSyncDataHandler(syncDataUC *usecases.SyncDataUseCase) *SyncDataHandler {
	return &SyncDataHandler{syncDataUC: syncDataUC}
}

// Handle implements the handlers.Handler interface for SyncDataHandler.
func (h *SyncDataHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := h.syncDataUC.GetAllDataForSync(ctx)
	if err != nil {
		log.Printf("SyncDataHandler: Error getting sync data: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, data)
}