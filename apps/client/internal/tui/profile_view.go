package tui

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProfileView wraps the components for the user profile screen.
type ProfileView struct {
	grid      *tview.Grid
	form      *tview.Form
	idField   *tview.InputField
	nameField *tview.InputField
}

// newProfileView creates and configures the profile view.
func newProfileView(actionsChan chan<- ui.UserAction) *ProfileView {
	pv := &ProfileView{}

	pv.form = tview.NewForm().SetItemPadding(1)

	pv.nameField = tview.NewInputField().SetLabel("Username:").SetFieldWidth(40)
	pv.idField = tview.NewInputField().SetLabel("Your ID:").SetFieldWidth(40)

	pv.form.AddFormItem(pv.nameField)
	pv.form.AddFormItem(pv.idField)
	pv.form.AddButton("Back to Chats", func() {
		actionsChan <- ui.SwitchChatAction{ChatID: ""}
	})
	pv.form.SetBorder(true).SetTitle("Your Profile")

	// --- Input Handlers ---

	// Capture input to make fields read-only, allowing only navigation.
	pv.nameField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Allow navigation keys
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			return event
		}
		// Handle 'y' for copy
		if event.Rune() == 'y' {
			clipboard.WriteAll(pv.idField.GetText())
			return nil
		}
		// Handle Escape to go back
		if event.Key() == tcell.KeyEscape {
			actionsChan <- ui.SwitchChatAction{ChatID: ""}
			return nil
		}
		// Block all other input
		return nil
	})

	pv.idField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Allow navigation keys
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			return event
		}
		// Handle 'y' for copy
		if event.Rune() == 'y' {
			clipboard.WriteAll(pv.idField.GetText())
			return nil
		}
		// Handle Escape to go back
		if event.Key() == tcell.KeyEscape {
			actionsChan <- ui.SwitchChatAction{ChatID: ""}
			return nil
		}
		// Block all other input
		return nil
	})

	// Use a grid to center the form
	pv.grid = tview.NewGrid().
		SetRows(0, 15, 0).
		SetColumns(0, 60, 0).
		AddItem(pv.form, 1, 1, 1, 1, 0, 0, true)

	return pv
}

// renderProfileView updates the view with the given model.
func (t *TUI) renderProfileView(model models.ProfileViewModel) {
	t.profileView.nameField.SetText(model.Username)
	t.profileView.idField.SetText(model.ID)
	t.profileView.form.SetTitle("Your Profile")
}
