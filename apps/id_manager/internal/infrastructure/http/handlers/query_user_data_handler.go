package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eaglechat/apps/id_manager/internal/application/usecases"
)

type QueryUserDataHandler struct {
	useCase usecases.UseCase[*usecases.QueryUserRequest, usecases.QueryUserResponse]
}

func NewQueryUserDataHandler(uc usecases.UseCase[*usecases.QueryUserRequest, usecases.QueryUserResponse]) *QueryUserDataHandler {
	return &QueryUserDataHandler{useCase: uc}
}

// Handle processes the user data query request.
func (h *QueryUserDataHandler) Handle(c *gin.Context) {
	var req usecases.QueryUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.useCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
