package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	idsStr := c.Query("ids")
	if idsStr != "" {
		if err := json.Unmarshal([]byte(idsStr), &req.Ids); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ids parameter"})
			return
		}
	}

	omitDisconnectedStr := c.DefaultQuery("omit_disconnected", "false")
	omitDisconnected, err := strconv.ParseBool(omitDisconnectedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid omit_disconnected parameter"})
		return
	}
	req.OmitDisconnected = omitDisconnected

	resp, err := h.useCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
