package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

// Message implements domain.Middleware.
func (m *Middleware) Message(targetID entities.UserID, message string) error {
	data, err := m.getUserData(targetID)
	if err != nil {
		log.Printf("User %s not found, storing message as pending.", targetID)
		return m.storeAsPending(targetID, message)
	}

	var sendErr error
	for i := range maxRetries {
		sendErr = m.p2pConnPool.Message(data.IP.String(), fmt.Sprint(m.ownPort), []byte(message))
		if sendErr == nil {
			return nil // Message sent successfully
		}
		log.Printf("Failed to send message to %s (attempt %d/%d): %v", targetID, i+1, maxRetries, sendErr)
		time.Sleep(retryDelay)
	}

	log.Printf("Failed to send message to %s after %d retries, storing as pending.", targetID, maxRetries)
	return m.storeAsPending(targetID, message)
}

// storeAsPending saves a message to the cache and notifies the ID manager.
func (m *Middleware) storeAsPending(targetID entities.UserID, message string) error {
	msgID := uuid.New().String()
	messageTarget := middleware_entities.NewMessageTarget(msgID, targetID)
	pendingMsg := middleware_entities.NewPendingMessage(messageTarget, []byte(message))

	// Store the message, ensuring it is not lost.
	if err := m.messageCache.StoreImmune(pendingMsg); err != nil {
		// This is a hard failure; we could not even store the message.
		return fmt.Errorf("failed to send message and failed to store it as pending: %w", err)
	}

	// Best-effort notification to the ID manager.
	idManager, err := m.iDManagerPool.GetAny()
	if err != nil {
		log.Printf("Failed to get ID manager to notify about pending message %s: %v", msgID, err)
		return nil // Return nil as the message is safely stored.
	}

	if err := idManager.NotifyOfPendingMessages([]middleware_entities.MessageTarget{messageTarget}); err != nil {
		log.Printf("Failed to notify ID manager about pending message %s: %v", msgID, err)
		// Return nil as the message is safely stored.
	}

	return nil
}
