package tui

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
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
func newChatView(actionsChan chan<- ui.UserAction) *ChatView {
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

	cv.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := cv.inputField.GetText()
			if text != "" {
				actionsChan <- ui.SendMessageAction{Content: text}
				cv.inputField.SetText("")
			}
		}
	})

	cv.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'p' {
			actionsChan <- ui.SwitchToProfileAction{}
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
