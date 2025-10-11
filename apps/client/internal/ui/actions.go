package ui

// UserAction is an interface for events initiated by the user in the UI.
type UserAction interface{ isUserAction() }

// SubmitRegistrationAction is triggered when the user submits their name for registration.
type SubmitRegistrationAction struct{ Name string }

func (SubmitRegistrationAction) isUserAction() {}

// SendMessageAction is triggered when the user sends a message to the active chat.
type SendMessageAction struct{ Content string }

func (SendMessageAction) isUserAction() {}

// SwitchChatAction is triggered when the user switches the active chat.
type SwitchChatAction struct{ ChatID string }

func (SwitchChatAction) isUserAction() {}
