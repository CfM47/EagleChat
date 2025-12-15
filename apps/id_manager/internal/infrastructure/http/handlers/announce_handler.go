package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnnounceHandler struct {
	useCase *usecases.AnnounceUseCase
}

func NewAnnounceHandler(uc *usecases.AnnounceUseCase) *AnnounceHandler {
	return &AnnounceHandler{useCase: uc}
}

func (h *AnnounceHandler) Handle(c *gin.Context) {
	// Get client ID from header
	clientID := c.GetHeader("X-Client-ID")
	if clientID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "X-Client-ID header is required",
		})
		return
	}

	//  Get client IP
	clientIPStr := c.ClientIP()
	ip := net.ParseIP(clientIPStr)
	if ip == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid client IP",
		})
		return
	}

	// Call use case
	if err := h.useCase.Execute(c.Request.Context(), clientID, ip); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
