package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterUserHandler struct {
	useCase usecases.UseCase[*usecases.RegisterUserRequest, *usecases.RegisterUserResponse]
}

func NewRegisterUserHandler(uc usecases.UseCase[*usecases.RegisterUserRequest, *usecases.RegisterUserResponse]) *RegisterUserHandler {
	return &RegisterUserHandler{useCase: uc}
}

func (h *RegisterUserHandler) Handle(c *gin.Context) {
	log.Println("Handling user registration request")
	var req usecases.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get IP from the request and add it to the request struct
	req.IP = c.ClientIP()
	log.Printf("New user registered with Ip: %s", req.IP)

	resp, err := h.useCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
