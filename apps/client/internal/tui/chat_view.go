package tui

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ChatView wraps the components for the main chat screen.
type ChatView struct {
	grid        *tview.Flex
	chatList    *tview.List
	messageView *tview.TextView
	inputField  *tview.InputField
}

// newChatView creates and configures the main chat view.
func newChatView(app *tview.Application, actionsChan chan<- ui.UserAction) *ChatView {
	ctx := ezlog.NewLoggerContext("chat-view-new")
	cv := &ChatView{}

	cv.chatList = tview.NewList().ShowSecondaryText(false)
	cv.chatList.SetBorder(true).SetTitle("Chats")

	cv.messageView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			cv.messageView.ScrollToEnd()
		})
	cv.messageView.SetBorder(true).SetTitle("Messages")

	cv.inputField = tview.NewInputField().
		SetLabel("Message: ").
		SetFieldWidth(0)

	// Main chat layout
	chatFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(cv.messageView, 0, 1, false).
		AddItem(cv.inputField, 3, 1, true) // Input field takes focus

	cv.grid = tview.NewFlex().
		AddItem(cv.chatList, 30, 1, false).
		AddItem(chatFlex, 0, 2, true)

	// --- Input Handlers ---

	// Handle j/k for vim-like navigation in the chat list.
	cv.chatList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			ezlog.Log(ctx).Debug("Chat list: 'j' pressed, moving down.")
			currentItem := cv.chatList.GetCurrentItem()
			itemCount := cv.chatList.GetItemCount()
			if itemCount > 0 {
				cv.chatList.SetCurrentItem((currentItem + 1) % itemCount)
			}
			return nil
		case 'k':
			ezlog.Log(ctx).Debug("Chat list: 'k' pressed, moving up.")
			currentItem := cv.chatList.GetCurrentItem()
			if currentItem > 0 {
				cv.chatList.SetCurrentItem(currentItem - 1)
			}
			return nil
		}
		return event
	})

	cv.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := cv.inputField.GetText()
			if text != "" {
				ezlog.Log(ctx).Infof("Sending message: %s", text)
				actionsChan <- ui.SendMessageAction{Content: text}
				cv.inputField.SetText("")
			}
		}
	})

	// Global input capture for the entire chat view
	cv.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If the user is typing, ignore all global shortcuts
		if cv.inputField.HasFocus() {
			// ...except for Tab, which we need to handle to escape the focus trap.
			if event.Key() == tcell.KeyTab {
				ezlog.Log(ctx).Debug("Tab pressed in input field, focusing chat list.")
				app.SetFocus(cv.chatList)
				return nil // Absorb the event
			}
			return event // Let the input field handle the key
		}

		// If the chat list has focus, Tab should switch to the input field
		if cv.chatList.HasFocus() && event.Key() == tcell.KeyTab {
			ezlog.Log(ctx).Debug("Tab pressed in chat list, focusing input field.")
			app.SetFocus(cv.inputField)
			return nil // Absorb the event
		}

		// Global shortcuts (only active when not typing)
		switch event.Rune() {
		case 'p':
			ezlog.Log(ctx).Debug("'p' pressed, switching to profile.")
			actionsChan <- ui.SwitchToProfileAction{}
			return nil
		case 'n':
			ezlog.Log(ctx).Debug("'n' pressed, switching to new contact view.")
			actionsChan <- ui.SwitchToNewContactViewAction{}
			return nil
		}
		return event
	})

	return cv
}

// renderChatView updates the view with the given model.
func (t *TUI) renderChatView(model models.ChatViewModel) {
	// Update Chat List
	t.chatView.chatList.Clear()
	for _, chat := range model.InactiveChats {
		// Capture chat for the closure
		chatCopy := chat
		var title string
		if chat.UnreadCount > 0 {
			title = fmt.Sprintf("%s (%d)", chat.Name, chat.UnreadCount)
		} else {
			title = chat.Name
		}
		var lastMessage string
		if chat.LastMessage == nil {
			lastMessage = "[No messages]"
		} else {
			lastMessage = *chat.LastMessage
		}

		// TODO: Timestamp: may be 0

		t.chatView.chatList.AddItem(title, lastMessage, 0, func() {
			t.actionsChan <- ui.SwitchChatAction{ChatID: chatCopy.ID}
		})
	}

	// Update Message View
	t.chatView.messageView.Clear()
	if model.ActiveChat != nil {
		t.chatView.messageView.SetTitle(fmt.Sprintf("Chat with %s", model.ActiveChat.Name))
		for _, msg := range model.ActiveChat.Messages {
			timestamp := msg.Timestamp.Format(time.Kitchen)
			fmt.Fprintf(t.chatView.messageView, "[gray]%s [white]%s: [white]%s\n", timestamp, msg.Author, msg.Content)
		}
	} else {
		t.chatView.messageView.SetTitle("Select a chat")
		fmt.Fprint(t.chatView.messageView, "Welcome to EagleChat!")
	}
}
