package controller

import (
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"log"
)

// Run starts the main application event loop.
func (c *Controller) Run() error {
	// 1. Check for existing user profile on startup.
	profile, err := c.userRepo.GetOwnProfile()

	if err != nil {
		// Assuming any error means no profile exists.
		log.Println("No profile found, entering registration.")
		c.appState = models.RegisterState
	} else {
		log.Println("Profile found, connecting...")
		c.appState = models.ChatState
		c.connectToMiddleware(profile)
	}

	// 2. Perform initial render.
	c.render()

	// 3. Start the main event loop.
	log.Println("Starting main event loop.")
	for {
		select {
		case action := <-c.ui.UserActions():
			c.handleUIAction(action)
		case msg := <-c.messageChannel:
			c.handleIncomingMessage(msg)
		}
	}
}

// handleUIAction dispatches actions from the UI to the appropriate handler.
func (c *Controller) handleUIAction(action ui.UserAction) {
	switch act := action.(type) {
	case ui.SubmitRegistrationAction:
		if c.appState == models.RegisterState {
			c.handleRegistration(act.Name)
		}
	case ui.SwitchChatAction:
		if c.appState == models.ChatState {
			c.handleSwitchChat(act.ChatID)
		}
	case ui.SendMessageAction:
		if c.appState == models.ChatState {
			c.handleSendMessage(act.Content)
		}
	}
}
