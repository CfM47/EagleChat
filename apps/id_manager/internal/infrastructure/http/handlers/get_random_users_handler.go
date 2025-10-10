package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"eaglechat/apps/id_manager/internal/application/usecases"
)

type GetRandomUsersHandler struct {
	useCase *usecases.GetRandomUsersUseCase
}

func NewGetRandomUsersHandler(uc *usecases.GetRandomUsersUseCase) *GetRandomUsersHandler {
	return &GetRandomUsersHandler{useCase: uc}
}

func (h *GetRandomUsersHandler) Handle(c *gin.Context) {
	amountStr := c.DefaultQuery("amount", "10") // Default to 10
	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount parameter"})
		return
	}

	req := usecases.GetRandomUsersRequest{
		Amount: amount,
	}

	resp, err := h.useCase.Execute(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
