package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"log"
)

func (c *Controller) handleSwitchChat(chatID string) {
	if c.activeChatID == chatID {
		return
	}

	log.Printf("Switching active chat to %s", chatID)
	c.activeChatID = chatID

	// Reset the unread count for the newly active chat.
	if err := c.repository.ResetUnreadCount(entities.UserID(chatID)); err != nil {
		log.Printf("Failed to reset unread count for chat %s: %v", chatID, err)
	}

	// Re-render the UI to show the new active chat.
	c.render()
}
