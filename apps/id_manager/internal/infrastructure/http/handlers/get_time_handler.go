package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetTimeHandler struct {
	useCase *usecases.GetTimeUseCase
}

func NewGetTimeHandler(uc *usecases.GetTimeUseCase) *GetTimeHandler {
	return &GetTimeHandler{useCase: uc}
}

func (h *GetTimeHandler) Handle(c *gin.Context) {
	resp := h.useCase.Execute(c.Request.Context(), struct{}{})

	encodedTime, err := resp.CurrentTime.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, encodedTime)
}
