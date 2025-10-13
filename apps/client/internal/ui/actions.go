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
// if ChatID is "", the chat view will be entered, this is used when switching back from profile view.
type SwitchChatAction struct{ ChatID string }

func (SwitchChatAction) isUserAction() {}

// SwitchToProfileAction is triggered when the user navigates to their profile view.
type SwitchToProfileAction struct{}

func (SwitchToProfileAction) isUserAction() {}

// SwitchToNewContactViewAction is triggered when the user navigates to the new contact screen.
type SwitchToNewContactViewAction struct{}

func (SwitchToNewContactViewAction) isUserAction() {}

// StartChatWithUserAction is triggered when the user submits a new user ID to start a chat.
type StartChatWithUserAction struct{ UserID string }

func (StartChatWithUserAction) isUserAction() {}
