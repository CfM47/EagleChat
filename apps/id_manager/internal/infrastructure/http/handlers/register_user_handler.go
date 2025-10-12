package handlers

import (
	"net"
	"net/http"

	"eaglechat/apps/id_manager/internal/application/usecases"
	"github.com/gin-gonic/gin"
)

type RegisterUserHandler struct {
	useCase *usecases.RegisterUserUseCase
}

func NewRegisterUserHandler(uc *usecases.RegisterUserUseCase) *RegisterUserHandler {
	return &RegisterUserHandler{useCase: uc}
}

func (h *RegisterUserHandler) Handle(c *gin.Context) {
	var req usecases.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ip := net.ParseIP(c.ClientIP())

	resp, err := h.useCase.Execute(c.Request.Context(), &req, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
