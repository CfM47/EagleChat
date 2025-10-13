package tui

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"

	"github.com/rivo/tview"
)

// TUI is a tview-based implementation of the ui.UI interface.
type TUI struct {
	app   *tview.Application
	pages *tview.Pages

	// Child views
	chatView       *ChatView
	registerView   *RegisterView
	profileView    *ProfileView
	newContactView *NewContactView

	actionsChan chan ui.UserAction
}

// New creates and initializes all UI components.
func New() *TUI {
	ctx := ezlog.NewLoggerContext("tui-new")
	ezlog.Log(ctx).Info("Initializing new TUI...")

	actionsChan := make(chan ui.UserAction)

	t := &TUI{
		app:         tview.NewApplication(),
		pages:       tview.NewPages(),
		actionsChan: actionsChan,
	}

	// Initialize child views
	ezlog.Log(ctx).Info("Creating child views...")
	t.chatView = newChatView(t.app, actionsChan)
	t.registerView = newRegisterView(actionsChan)
	t.profileView = newProfileView(actionsChan)
	t.newContactView = newNewContactView(actionsChan)
	ezlog.Log(ctx).Info("Child views created.")

	// Add pages
	ezlog.Log(ctx).Info("Adding pages...")
	t.pages.AddPage("register", t.registerView.grid, true, false)
	t.pages.AddPage("chat", t.chatView.grid, true, false)
	t.pages.AddPage("profile", t.profileView.grid, true, false)
	t.pages.AddPage("new_contact", t.newContactView.grid, true, false)
	ezlog.Log(ctx).Info("Pages added.")

	t.app.SetRoot(t.pages, true).EnableMouse(true)
	ezlog.Log(ctx).Info("TUI initialization complete.")
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
	ctx := ezlog.NewLoggerContext("tui-render")
	t.app.QueueUpdateDraw(func() {
		ezlog.Log(ctx).Infof("Rendering new state: %v", model.CurrentState)
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
		case models.NewContactState:
			t.renderNewContactView(ctx, model.NewContactView)
			t.pages.SwitchToPage("new_contact")
		}
	})
}
