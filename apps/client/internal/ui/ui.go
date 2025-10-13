package ui

import "eaglechat/apps/client/internal/ui/models"

// UI defines the contract for the user interface.
type UI interface {
	// Render provides the UI with the complete model it needs to draw itself.
	// The UI should be stateless and simply reflect the given model.
	Render(model models.UIModel)

	// UserActions returns a channel that the application logic can listen to
	// for actions initiated by the user.
	UserActions() <-chan UserAction

	// Run starts the UI and blocks until it's terminated.
	Run() error
}
