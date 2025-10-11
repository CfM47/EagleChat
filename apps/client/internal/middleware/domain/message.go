package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto"
	"encoding/json"
	"fmt"
	"log"
	"net"
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

	ip, err := m.getUserIfConnected(target.ID)
	if err != nil {
		return m.storeAsPending(pendingMsg)
	}

	if err := m.sendP2PMessage(ip, msgBytes); err != nil {
		log.Printf("failed to send message to %s, storing as pending: %v", target.ID, err)
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
func (m *Middleware) sendP2PMessage(targetIP net.IP, msgBytes []byte) error {
	var lastErr error
	for range maxRetries {
		err := m.p2pConnPool.Message(targetIP.String(), fmt.Sprint(m.ownPort), msgBytes)
		if err == nil {
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
	return nil
}

func (m *Middleware) getUserIfConnected(userID entities.UserID) (net.IP, error) {
	data, err := m.getUserData([]entities.UserID{userID}, true)
	if err != nil {
		return nil, err
	}

	if len(data) != 1 {
		return nil, fmt.Errorf("user not found")
	}

	return *data[userID].IP, nil
}
