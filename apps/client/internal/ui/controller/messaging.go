package controller

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
)

func (c *Controller) handleSendMessage(ctx context.Context, content string) {
	if c.connectedMiddleware == nil {
		ezlog.Log(ctx).Warnf("Cannot send message: middleware not connected.")
		return
	}
	if c.activeChatID == "" {
		ezlog.Log(ctx).Warnf("Cannot send message: no active chat.")
		return
	}

	// Get the full user entity for the target.
	targetUser, err := c.repository.GetUser(entities.UserID(c.activeChatID))
	if err != nil {
		ezlog.Log(ctx).Errorf("Cannot send message: could not find user %s: %v", c.activeChatID, err)
		return
	}

	self, err := c.repository.GetOwnProfile()
	if err != nil {
		ezlog.Log(ctx).Errorf("Could not get own profile: %v", err)
		return
	}

	message := entities.NewMessage(self.User, targetUser, content)

	c.repository.SaveMessage(message)

	ezlog.Log(ctx).Infof("Sending message to %s", c.activeChatID)
	go func() {
		if err := c.connectedMiddleware.Message(ctx, targetUser, message); err != nil {
			ezlog.Log(ctx).Errorf("Failed to send message to %s: %v", c.activeChatID, err)
		}
	}()

	c.render()
}

func (c *Controller) handleIncomingMessage(ctx context.Context, msg entities.Message) {
	ezlog.Log(ctx).Infof("Received message from %s", msg.Sender.ID)

	if err := c.repository.SaveMessage(msg); err != nil {
		ezlog.Log(ctx).Errorf("Failed to save incoming message: %v", err)
		return
	}

	// If the message is not for the currently active chat, increment unread count.
	if string(msg.Sender.ID) != c.activeChatID {
		if err := c.repository.IncrementUnreadCount(msg.Sender.ID); err != nil {
			ezlog.Log(ctx).Errorf("Failed to increment unread count: %v", err)
		}
	}

	// Re-render the UI to show the new message/unread count.
	c.render()
}
