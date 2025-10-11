package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto"
	"encoding/json"
	"fmt"
	"log"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"

	"github.com/google/uuid"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

// Message is the high-level orchestrator for sending a message.
func (m *Middleware) Message(target entities.User, message entities.Message) error {
	log.SetPrefix("[Message] ")

	pendingMsg, msgBytes, err := m.composeP2PMessage(target, message)
	if err != nil {
		return fmt.Errorf("failed to compose message: %w", err)
	}

	if err := m.sendP2PMessage(target, msgBytes); err != nil {
		log.Printf("Failed to send message to %s, storing as pending: %v", target.ID, err)
		return m.storeAsPending(pendingMsg)
	}

	return nil
}

// composeP2PMessage handles the creation and encryption of a message.
func (m *Middleware) composeP2PMessage(target entities.User, message entities.Message) (middleware_entities.PendingMessage, []byte, error) {
	// 1. Marshal the domain-level message object.
	innerMsgBytes, err := json.Marshal(message)
	if err != nil {
		return middleware_entities.PendingMessage{}, nil, fmt.Errorf("failed to marshal inner message: %w", err)
	}

	// 2. Encrypt the message into a secure envelope.
	envelope, err := simplecrypto.Seal(innerMsgBytes, &m.sk, &target.PublicKey)
	if err != nil {
		return middleware_entities.PendingMessage{}, nil, fmt.Errorf("failed to seal message: %w", err)
	}
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return middleware_entities.PendingMessage{}, nil, fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// 3. Create the P2P wire message (PendingMessage).
	msgTarget := middleware_entities.NewMessageTarget(uuid.New().String(), target.ID)
	pendingMsg := middleware_entities.NewPendingMessage(msgTarget, envelopeBytes)
	pendingMsgBytes, err := json.Marshal(pendingMsg)
	if err != nil {
		return middleware_entities.PendingMessage{}, nil, fmt.Errorf("failed to marshal pending message: %w", err)
	}

	return pendingMsg, pendingMsgBytes, nil
}

// sendP2PMessage handles the network logic of sending a message with retries.
func (m *Middleware) sendP2PMessage(target entities.User, msgBytes []byte) error {
	data, err := m.getUserData(target.ID)
	if err != nil {
		return fmt.Errorf("user %s not found: %w", target.ID, err)
	}

	var lastErr error
	for range maxRetries {
		if err := m.p2pConnPool.Message(data.IP.String(), fmt.Sprint(m.ownPort), msgBytes); err == nil {
			return nil // Message sent successfully
		}
		lastErr = err
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to send message after %d retries: %w", maxRetries, lastErr)
}

// storeAsPending saves a message to the cache.
func (m *Middleware) storeAsPending(pendingMsg middleware_entities.PendingMessage) error {
	if err := m.messageCache.StoreImmune(pendingMsg); err != nil {
		return fmt.Errorf("failed to store pending message: %w", err)
	}
	// TODO: Notify ID Manager
	return nil
}
