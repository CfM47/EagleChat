package handlers

import (
	"context"
	"eaglechat/apps/id_manager/internal/application/ports"
	"eaglechat/common/ezlog"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NotifyUpdateHandler handles POST requests for /notify-update endpoint.
type NotifyUpdateHandler struct {
	notifier ports.Notifier
}

// NewNotifyUpdateHandler creates a new NotifyUpdateHandler.
func NewNotifyUpdateHandler(notifier ports.Notifier) *NotifyUpdateHandler {
	return &NotifyUpdateHandler{notifier: notifier}
}

type notifyUpdateRequest struct {
	SourceAddress string `json:"source_address"`
}

// Handle implements the handlers.Handler interface for NotifyUpdateHandler.
func (h *NotifyUpdateHandler) Handle(c *gin.Context) {

	logCtx := ezlog.NewLoggerContext("Notify Update Handler")

	var req notifyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.SourceAddress == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "source_address is required"})
		return
	}

	ezlog.Log(logCtx).Debugf("Received update-notification from source %s", req.SourceAddress)

	// We run this in a goroutine to avoid blocking the caller.
	// The caller (another ID manager) doesn't need to wait for our sync to complete.
	go func(source string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := h.notifier.TriggerSyncFromPeer(ctx, source); err != nil {
			// Log the error for observability, but we don't need to return an error to the caller.
			ezlog.Log(logCtx).Errorf("NotifyUpdateHandler: failed to trigger sync from peer %s: %v", req.SourceAddress, err)
		}
	}(fmt.Sprintf("%s:%s", req.SourceAddress, h.notifier.GetPort()))

	c.JSON(http.StatusOK, gin.H{"status": "sync triggered"})
}
