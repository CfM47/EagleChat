package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/ui/models"
	"log"
)

// render queries the repositories and builds the UI model for the current state.
func (c *Controller) render() {
	model := models.UIModel{CurrentState: c.appState}

	switch c.appState {
	case models.RegisterState:
		model.RegisterView = models.RegisterViewModel{
			PromptMessage: "Welcome to EagleChat! Please enter your name to register.",
			ErrorMessage:  c.error,
		}

	case models.ChatState:
		overviews, err := c.messageRepo.GetAllChatOverviews()
		if err != nil {
			log.Printf("Failed to get chat overviews: %v", err)
		}
		activeChatMessages, err := c.messageRepo.GetChat(entities.UserID(c.activeChatID))
		if err != nil {
			log.Printf("Failed to get active chat: %v", err)
		}

		inactiveChats := make([]models.InactiveChatModel, len(overviews))
		for i, ov := range overviews {
			inactiveChats[i] = models.InactiveChatModel{
				ID:          string(ov.Partner.ID),
				Name:        ov.Partner.Name,
				LastMessage: ov.LastMessage,
				UnreadCount: ov.UnreadCount,
			}
		}

		var activeChat *models.ActiveChatModel
		if c.activeChatID != "" {
			msgs := make([]models.MessageModel, len(activeChatMessages))
			for i, msg := range activeChatMessages {
				msgs[i] = models.MessageModel{
					Timestamp: msg.CreatedTime,
					Author:    msg.Sender.Name,
					Content:   msg.Content,
				}
			}
			activeChat = &models.ActiveChatModel{
				ID:       c.activeChatID,
				Messages: msgs,
			}
		}

		model.ChatView = models.ChatViewModel{
			ActiveChat:    activeChat,
			InactiveChats: inactiveChats,
		}
	}

	log.Printf("Rendering state: %v", c.appState)
	c.ui.Render(model)
}
