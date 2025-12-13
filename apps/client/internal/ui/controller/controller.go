package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"
	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/simplecrypto/rsa"
)

// Controller orchestrates the application logic, acting as a bridge between the
// domain (Middleware) and the UI.
type Controller struct {
	ui         ui.UI
	registerer services.Registerer
	connector  services.Connector
	repository repositories.ClientRepository
	CAPubkey   rsa.PublicKey

	// Live instances, populated after connecting
	connectedMiddleware services.Middleware
	messageChannel      <-chan entities.Message

	// Internal state
	appState     models.AppState
	activeChatID string
	error        string
}

// New creates a new controller.
func New(
	ui ui.UI,
	registerer services.Registerer,
	connector services.Connector,
	repository repositories.ClientRepository,
	CAPubkey rsa.PublicKey,
) *Controller {
	return &Controller{
		ui:             ui,
		registerer:     registerer,
		connector:      connector,
		repository:     repository,
		messageChannel: nil,
		CAPubkey:       CAPubkey,
	}
}
