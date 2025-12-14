package handlers

import (
	"net/http"

	"eaglechat/apps/id_manager/internal/application/usecases"

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
	c.JSON(http.StatusOK, resp.CurrentTime)
}
