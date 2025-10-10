package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eaglechat/apps/id_manager/internal/application/usecases"
)

type AddPendingMessagesHandler struct {
	useCase *usecases.AddPendingMessagesUseCase
}

func NewAddPendingMessagesHandler(uc *usecases.AddPendingMessagesUseCase) *AddPendingMessagesHandler {
	return &AddPendingMessagesHandler{useCase: uc}
}

func (h *AddPendingMessagesHandler) Handle(c *gin.Context) {
	// TODO: Implement a robust way to infer client ID, for now, we'll use a header.
	cacherID := c.GetHeader("X-Client-ID")
	if cacherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var req usecases.AddPendingMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	req.CacherID = cacherID

	if err := h.useCase.Execute(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
