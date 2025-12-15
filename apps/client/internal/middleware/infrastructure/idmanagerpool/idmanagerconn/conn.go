package idmanagerconn

import (
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/common/simplecrypto/rsa"
)

// FIXME: expect all manager responses to be signed by the manager
type idManagerConnectionImpl struct {
	client    *http.Client
	baseURL   string
	ownID     entities.UserID
	managerPK rsa.PublicKey
}

func NewIDManagerConnectionImpl(client *http.Client, baseURL string, ownID entities.UserID, managerPK rsa.PublicKey) *idManagerConnectionImpl {
	return &idManagerConnectionImpl{
		client:    client,
		baseURL:   baseURL,
		ownID:     ownID,
		managerPK: managerPK,
	}
}

var _ managerpool_entities.IDManagerConnection = (*idManagerConnectionImpl)(nil)

func (c *idManagerConnectionImpl) BaseURL() string {
	return c.baseURL
}
