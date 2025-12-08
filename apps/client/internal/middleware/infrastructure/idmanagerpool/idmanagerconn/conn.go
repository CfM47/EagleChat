package idmanagerconn

import (
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
)

type idManagerConnectionImpl struct {
	client  *http.Client
	baseURL string
	ownID   entities.UserID
}

var _ managerpool_entities.IDManagerConnection = (*idManagerConnectionImpl)(nil)

func (c *idManagerConnectionImpl) BaseURL() string {
	return c.baseURL
}
