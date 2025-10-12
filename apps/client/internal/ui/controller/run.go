package controller

import (
	"context"
	"eaglechat/apps/client/internal/ui"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"
	"log"
)

// Run starts the main application event loop.
func (c *Controller) Run() error {
	ctx := ezlog.NewLoggerContext("controller")

	go func() {
		if err := c.ui.Run(); err != nil {
			ezlog.Log(ctx).Fatalf("UI exited with error: %v", err)
		}
	}()

	ezlog.Log(ctx).Info("UI started")

	ezlog.Log(ctx).Info("Checking for existing profile...")
	profile, err := c.repository.GetOwnProfile()

	if err != nil {
		ezlog.Log(ctx).Info("No profile found, entering registration.")
		c.appState = models.RegisterState
	} else {
		ezlog.Log(ctx).Info("Profile found, connecting...")
		c.appState = models.ChatState

		connectorCtx := ezlog.WithComponentPrefix(ctx, "middleware-connector")
		c.connectToMiddleware(connectorCtx, profile)
	}

	// 2. Perform initial render.
	c.render()

	// 3. Start the main event loop.
	log.Println("Starting main event loop.")
	for {
		select {
		case action := <-c.ui.UserActions():
			c.handleUIAction(ctx, action)
		case msg := <-c.messageChannel:
			messageCtx := ezlog.WithComponentPrefix(ctx, "incoming-message")
			c.handleIncomingMessage(messageCtx, msg)
		}
	}
}

// handleUIAction dispatches actions from the UI to the appropriate handler.
func (c *Controller) handleUIAction(ctx context.Context, action ui.UserAction) {
	switch act := action.(type) {
	case ui.SubmitRegistrationAction:
		if c.appState == models.RegisterState {
			ctx = ezlog.WithComponentPrefix(ctx, "registration")
			c.handleRegistration(ctx, act.Name)
		}
	case ui.SwitchChatAction:
		if c.appState == models.ChatState {
			ctx = ezlog.WithComponentPrefix(ctx, "switch-chat")
			c.handleSwitchChat(ctx, act.ChatID)
		}
	case ui.SendMessageAction:
		if c.appState == models.ChatState {
			ctx = ezlog.WithComponentPrefix(ctx, "send-message")
			c.handleSendMessage(ctx, act.Content)
		}
	}
}
