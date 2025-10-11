package tui

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"fmt"

	"github.com/rivo/tview"
)

// RegisterView wraps the components for the registration screen.
type RegisterView struct {
	grid *tview.Grid
	form *tview.Form
}

// newRegisterView creates and configures the registration view.
func newRegisterView(actionsChan chan<- ui.UserAction) *RegisterView {
	form := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1)

	form.AddInputField("Name", "", 20, nil, nil).
		AddButton("Register", func() {
			name := form.GetFormItem(0).(*tview.InputField).GetText()
			if name != "" {
				actionsChan <- ui.SubmitRegistrationAction{Name: name}
			}
		}).
		SetBorder(true).SetTitle("Enter Your Name")

	// Use a grid to center the form
	grid := tview.NewGrid().
		SetRows(0, 15, 0).
		SetColumns(0, 60, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	return &RegisterView{
		grid: grid,
		form: form,
	}
}

func (t *TUI) renderRegisterView(model models.RegisterViewModel) {
	if model.ErrorMessage != "" {
		title := fmt.Sprintf("[red]%s", model.ErrorMessage)
		t.registerView.form.SetTitle(title)
	} else {
		t.registerView.form.SetTitle(model.PromptMessage)
	}
}
