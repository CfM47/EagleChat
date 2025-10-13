package controller

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
)

func (c *Controller) handleStartChatWithUser(ctx context.Context, userID string) {
	ezlog.Log(ctx).Infof("Attempting to start chat with user: %s", userID)

	// 1. Query the middleware for the user
	user, err := c.connectedMiddleware.QueryUser(entities.UserID(userID))
	if err != nil {
		ezlog.Log(ctx).Warnf("Failed to query user %s: %v", userID, err)
		c.setError("User not found or an error occurred.")
		c.render()
		return
	}

	// 2. Save the user to the local repository
	if err := c.repository.SaveUser(user); err != nil {
		ezlog.Log(ctx).Errorf("Failed to save new user %s: %v", userID, err)
		c.setError("Failed to save new contact.")
		c.render()
		return
	}

	ezlog.Log(ctx).Infof("Successfully found and saved user %s. Switching to chat.", user.Name)

	// 3. Switch to the new chat
	c.handleSwitchChat(ctx, string(user.ID))
}
