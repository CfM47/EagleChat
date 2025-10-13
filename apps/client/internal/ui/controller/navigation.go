package controller

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"
)

func (c *Controller) handleSwitchChat(ctx context.Context, chatID string) {
	if c.activeChatID != "" && c.activeChatID == chatID {
		ezlog.Log(ctx).Debug("Chat is already active, no switch needed")
		return
	}

	c.activeChatID = chatID

	c.appState = models.ChatState

	// Reset the unread count for the newly active chat.
	if err := c.repository.ResetUnreadCount(entities.UserID(chatID)); err != nil {
		ezlog.Log(ctx).Errorf("Failed to reset unread count for chat %s: %v", chatID, err)
	}

	// Re-render the UI to show the new active chat.
	c.render()
}

func (c *Controller) handleSwitchToProfile(ctx context.Context) {
	ezlog.Log(ctx).Info("Switching to profile view")

	c.appState = models.ProfileState

	c.render()
}

func (c *Controller) handleSwitchToNewContactView(ctx context.Context) {
	ezlog.Log(ctx).Info("Switching to new contact view")
	c.appState = models.NewContactState
	c.render()
}
