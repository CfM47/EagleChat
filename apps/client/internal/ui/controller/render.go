package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"
	"log"
)

// render queries the repositories and builds the UI model for the current state.
func (c *Controller) render() {
	ctx := ezlog.NewLoggerContext("render")

	model := models.UIModel{CurrentState: c.appState}

	switch c.appState {
	case models.RegisterState:
		model.RegisterView = models.RegisterViewModel{
			PromptMessage: "Welcome to EagleChat! Please enter your name to register.",
			ErrorMessage:  c.error,
		}
	case models.ProfileState:
		profile, err := c.repository.GetOwnProfile()
		if err != nil {
			ezlog.Log(ctx).Errorf("Failed to get own profile: %v", err)
			return
		}
		model.ProfileView = models.ProfileViewModel{
			ID:       string(profile.User.ID),
			Username: profile.User.Name,
		}
	case models.NewContactState:
		model.NewContactView = models.NewContactViewModel{
			ErrorMessage: c.error,
		}

	case models.ChatState:
		overviews, err := c.repository.GetAllChatOverviews()
		if err != nil {
			ezlog.Log(ctx).Errorf("Failed to get chat overviews: %v", err)
		}

		activeChatMessages, err := c.repository.GetChat(entities.UserID(c.activeChatID))
		ezlog.Log(ctx).Debugf("Active chat ID: %s, messages count: %d", c.activeChatID, len(activeChatMessages))

		if err != nil {
			ezlog.Log(ctx).Errorf("Failed to get active chat messages: %v", err)
		}
		var user entities.User
		if c.activeChatID != "" {
			user, err = c.repository.GetUser(entities.UserID(c.activeChatID))
			if err != nil {
				ezlog.Log(ctx).Errorf("Failed to get user for active chat ID %s: %v", c.activeChatID, err)
			}
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
				Name:     user.Name,
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
