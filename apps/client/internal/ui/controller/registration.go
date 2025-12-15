package controller

import (
	"context"
	"fmt"
	"log"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

func (c *Controller) handleRegistration(ctx context.Context, name string) {
	ezlog.Log(ctx).Infof("Attempting to register user with name: %s", name)

	// 1. Generate key pair
	privKey, _, err := rsa.GenerateKeyPair()
	if err != nil {
		msg := fmt.Sprintf("Failed to generate key pair: %v", err)
		log.Print(msg)
		c.appState = models.RegisterState
		c.setError(msg)
		return
	}

	// 2. Call the registerer to get a user ID from the ID Manager
	user, err := c.registerer.Register(ctx, name, *privKey, c.CAPubkey)
	if err != nil {
		msg := fmt.Sprintf("Failed to register with ID manager: %v", err)
		ezlog.Log(ctx).Error(msg)
		c.setError(msg)
		return
	}

	// 3. Save the full profile locally
	profile := entities.OwnProfile{
		User:       user,
		PrivateKey: *privKey,
	}
	if err := c.repository.SaveOwnProfile(profile); err != nil {
		ezlog.Log(ctx).Errorf("Failed to save profile: %v", err)
		return
	}

	ezlog.Log(ctx).Info("Registration successful. Connecting to middleware...")

	// 4. Connect to the middleware and transition to chat state
	c.appState = models.ChatState
	c.connectToMiddleware(ctx, profile)
	c.render()
}
