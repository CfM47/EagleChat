package controller

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
)

func (c *Controller) handleSwitchChat(ctx context.Context, chatID string) {
	if c.activeChatID == chatID {
		ezlog.Log(ctx).Debug("Chat is already active, no switch needed")
		return
	}

	ezlog.Log(ctx).Infof("Switching active chat to %s", chatID)
	c.activeChatID = chatID

	// Reset the unread count for the newly active chat.
	if err := c.repository.ResetUnreadCount(entities.UserID(chatID)); err != nil {
		ezlog.Log(ctx).Errorf("Failed to reset unread count for chat %s: %v", chatID, err)
	}

	// Re-render the UI to show the new active chat.
	c.render()
}
