package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eaglechat/apps/id_manager/internal/application/usecases"
)

type QueryPendingMessagesHandler struct {
	useCase *usecases.QueryPendingMessagesUseCase
}

func NewQueryPendingMessagesHandler(uc *usecases.QueryPendingMessagesUseCase) *QueryPendingMessagesHandler {
	return &QueryPendingMessagesHandler{useCase: uc}
}

func (h *QueryPendingMessagesHandler) Handle(c *gin.Context) {
	var targetID *string
	if tid, ok := c.GetQuery("target_id"); ok {
		targetID = &tid
	}

	var get_cachers bool
	if gc, ok := c.GetQuery("get_cachers"); ok {
		get_cachers = gc == "true"
	}

	req := usecases.QueryPendingMessagesRequest{
		TargetID:   targetID,
		GetCachers: get_cachers,
	}

	resp, err := h.useCase.Execute(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
