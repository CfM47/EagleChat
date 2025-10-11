package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/ui/models"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"fmt"
	"log"
)

func (c *Controller) handleRegistration(name string) {
	log.Printf("Attempting to register user with name: %s", name)

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
	user, err := c.registerer.Register(name, *privKey)
	if err != nil {
		msg := fmt.Sprintf("Failed to register with ID manager: %v", err)
		log.Print(msg)
		c.setError(msg)
		return
	}

	// 3. Save the full profile locally
	profile := entities.OwnProfile{
		User:       user,
		PrivateKey: *privKey,
	}
	if err := c.userRepo.SaveOwnProfile(profile); err != nil {
		log.Fatalf("Failed to save profile: %v", err)
		return
	}

	log.Println("Registration successful. Connecting to middleware...")

	// 4. Connect to the middleware and transition to chat state
	c.appState = models.ChatState
	c.connectToMiddleware(profile)
	c.render()
}
