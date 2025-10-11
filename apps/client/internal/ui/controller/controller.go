package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	message_repo "eaglechat/apps/client/internal/domain/repositories/message"
	user_repo "eaglechat/apps/client/internal/domain/repositories/user"
	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
)

// Controller orchestrates the application logic, acting as a bridge between the
// domain (Middleware) and the UI.
type Controller struct {
	ui          ui.UI
	registerer  services.Registerer
	connector   services.Connector
	userRepo    user_repo.UserRepository
	messageRepo message_repo.MessageRepository

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
	userRepo user_repo.UserRepository,
	messageRepo message_repo.MessageRepository,
) *Controller {
	return &Controller{
		ui:             ui,
		registerer:     registerer,
		connector:      connector,
		userRepo:       userRepo,
		messageRepo:    messageRepo,
		messageChannel: nil,
	}
}

