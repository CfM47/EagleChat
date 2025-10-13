package tui

import (
	"context"
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewContactView wraps the components for the "new contact" screen.
type NewContactView struct {
	grid *tview.Grid
	form *tview.Form
}

// newNewContactView creates and configures the new contact view.
func newNewContactView(actionsChan chan<- ui.UserAction) *NewContactView {
	ctx := ezlog.NewLoggerContext("tui-new-contact-view")

	ncv := &NewContactView{}

	ncv.form = tview.NewForm().
		SetItemPadding(1)

	ncv.form.AddInputField("User ID", "", 40, nil, nil).
		AddButton("Start Chat", func() {
			id := ncv.form.GetFormItem(0).(*tview.InputField).GetText()
			if id != "" {
				actionsChan <- ui.StartChatWithUserAction{UserID: id}
			}
		}).
		SetBorder(true).SetTitle("Start New Chat")

	// --- Input Handlers ---
	ncv.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ezlog.Log(ctx).Infof("NewContactView captured key: %v", event)
		if event.Key() == tcell.KeyEscape {
			actionsChan <- ui.SwitchChatAction{ChatID: ""}
			return nil
		}
		return event
	})

	// Use a grid to center the form
	ncv.grid = tview.NewGrid().
		SetRows(0, 15, 0).
		SetColumns(0, 60, 0).
		AddItem(ncv.form, 1, 1, 1, 1, 0, 0, true)

	return ncv
}

// renderNewContactView updates the view with the given model.
func (t *TUI) renderNewContactView(ctx context.Context, model models.NewContactViewModel) {
	ezlog.Log(ctx).Debugf("Rendering NewContactView with model: %+v", model)

	if model.ErrorMessage != "" {
		ezlog.Log(ctx).Infof("NewContactView error: %s", model.ErrorMessage)
		title := fmt.Sprintf("[red]Error: %s", model.ErrorMessage)
		t.newContactView.form.SetTitle(title)
	} else {
		t.newContactView.form.SetTitle("Start New Chat")
	}
}
