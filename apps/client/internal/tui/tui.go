package tui

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"

	"github.com/rivo/tview"
)

// TUI is a tview-based implementation of the ui.UI interface.
type TUI struct {
	app   *tview.Application
	pages *tview.Pages

	// Child views
	chatView     *ChatView
	registerView *RegisterView
	profileView  *ProfileView

	actionsChan chan ui.UserAction
}

// New creates and initializes all UI components.
func New() *TUI {
	actionsChan := make(chan ui.UserAction)

	t := &TUI{
		app:         tview.NewApplication(),
		pages:       tview.NewPages(),
		actionsChan: actionsChan,
	}

	// Initialize child views
	t.chatView = newChatView(actionsChan)
	t.registerView = newRegisterView(actionsChan)
	t.profileView = newProfileView(actionsChan)

	// Add pages
	t.pages.AddPage("register", t.registerView.grid, true, false)
	t.pages.AddPage("chat", t.chatView.grid, true, false)
	t.pages.AddPage("profile", t.profileView.grid, true, false)

	t.app.SetRoot(t.pages, true).EnableMouse(true)
	return t
}

// Run starts the UI application and blocks until it is stopped.
func (t *TUI) Run() error {
	return t.app.Run()
}

// UserActions returns the channel that the UI emits user actions on.
func (t *TUI) UserActions() <-chan ui.UserAction {
	return t.actionsChan
}

// Render updates the UI widgets to reflect the given model.
func (t *TUI) Render(model models.UIModel) {
	t.app.QueueUpdateDraw(func() {
		switch model.CurrentState {
		case models.RegisterState:
			t.renderRegisterView(model.RegisterView)
			t.pages.SwitchToPage("register")
		case models.ChatState:
			t.renderChatView(model.ChatView)
			t.pages.SwitchToPage("chat")
		case models.ProfileState:
			t.renderProfileView(model.ProfileView)
			t.pages.SwitchToPage("profile")
		}
	})
}