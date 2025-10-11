package models

import "time"

// AppState defines the possible high-level states of the application.
type AppState int

const (
	RegisterState AppState = iota
	ChatState
)

// MessageModel represents a single message for display in the UI.
type MessageModel struct {
	Timestamp time.Time
	Author    string
	Content   string
}

// InactiveChatModel represents a chat in the side panel list.
type InactiveChatModel struct {
	ID          string
	Name        string
	LastMessage string // For the preview text
	UnreadCount int
}

// ActiveChatModel represents the currently selected, visible chat.
type ActiveChatModel struct {
	ID       string
	Name     string
	Messages []MessageModel // The full message history to display
}

// ChatViewModel holds the data for the main chat screen.
type ChatViewModel struct {
	ActiveChat    *ActiveChatModel
	InactiveChats []InactiveChatModel
}

// RegisterViewModel holds the data needed to render the registration screen.
type RegisterViewModel struct {
	PromptMessage string
	ErrorMessage  string
}

// UIModel is the top-level model for the entire application.
// It clearly defines the current state and provides the specific model for that state.
type UIModel struct {
	CurrentState AppState
	RegisterView RegisterViewModel
	ChatView     ChatViewModel
}
