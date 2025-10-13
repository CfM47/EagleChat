package handlers

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"eaglechat/apps/id_manager/internal/application/usecases"
)

type QueryPendingMessagesHandler struct {
	useCase usecases.UseCase[*usecases.QueryPendingMessagesRequest, *usecases.QueryPendingMessagesResponse]
}

func NewQueryPendingMessagesHandler(uc usecases.UseCase[*usecases.QueryPendingMessagesRequest, *usecases.QueryPendingMessagesResponse]) *QueryPendingMessagesHandler {
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

	// TODO: Implement a robust way to infer client ID, for now, we'll use a header.
	querierID := c.GetHeader("X-Client-ID")

	req := usecases.QueryPendingMessagesRequest{
		TargetID:   targetID,
		GetCachers: get_cachers,
		QuerierID:  querierID,
		IP:         net.ParseIP(c.ClientIP()),
	}

	resp, err := h.useCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
