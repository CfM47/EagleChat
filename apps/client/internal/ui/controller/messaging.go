package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"log"
)

func (c *Controller) handleSendMessage(content string) {
	if c.connectedMiddleware == nil {
		log.Println("Cannot send message: middleware not connected.")
		return
	}
	if c.activeChatID == "" {
		log.Println("Cannot send message: no active chat.")
		return
	}

	// Get the full user entity for the target.
	targetUser, err := c.userRepo.Get(entities.UserID(c.activeChatID))
	if err != nil {
		log.Printf("Cannot send message: could not find user %s: %v", c.activeChatID, err)
		return
	}

	self, err := c.userRepo.GetOwnProfile()
	if err != nil {
		log.Printf("Could not get own profile: %v", err)
		return
	}

	message := entities.NewMessage(self.User, content)

	log.Printf("Sending message to %s", c.activeChatID)
	go func() {
		if err := c.connectedMiddleware.Message(targetUser, message); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}()
}

func (c *Controller) handleIncomingMessage(msg entities.Message) {
	log.Printf("Handling incoming message from %s", msg.Sender.ID)

	if err := c.messageRepo.Save(msg); err != nil {
		log.Printf("Failed to save incoming message: %v", err)
		return
	}

	// If the message is not for the currently active chat, increment unread count.
	if string(msg.Sender.ID) != c.activeChatID {
		if err := c.messageRepo.IncrementUnreadCount(msg.Sender.ID); err != nil {
			log.Printf("Failed to increment unread count: %v", err)
		}
	}

	// Re-render the UI to show the new message/unread count.
	c.render()
}
