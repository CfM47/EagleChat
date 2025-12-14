package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// Message is the high-level orchestrator for sending a message.
func (m *Middleware) Message(ctx context.Context, target entities.User, message entities.Message) error {
	pendingMsg, err := m.composeP2PMessage(target, message)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to compose pending message: %v", err)
		return fmt.Errorf("failed to compose pending message: %w", err)
	}

	ip, err := m.getUserIfConnected(target.ID)
	if err != nil {
		ezlog.Log(ctx).Warnf("User %s not connected, storing message as pending: %v", target.ID, err)
		return m.storeAsPending(pendingMsg)
	}

	targetData := middleware_entities.NewUserData(target, &ip)

	if err := m.clientConnPool.Message(ctx, []middleware_entities.PendingMessage{pendingMsg}, targetData); err != nil {
		ezlog.Log(ctx).Warnf("Failed to send message to %s, storing as pending: %v", target.ID, err)
		return m.storeAsPending(pendingMsg)
	}

	return nil
}

// composeP2PMessage handles the creation and encryption of a message.
func (m *Middleware) composeP2PMessage(target entities.User, message entities.Message) (middleware_entities.PendingMessage, error) {
	// 1. Marshal the domain-level message object.
	innerMsgBytes, err := json.Marshal(message)
	if err != nil {
		return middleware_entities.PendingMessage{}, fmt.Errorf("failed to marshal inner message: %w", err)
	}

	// 2. Encrypt the message into a secure envelope.
	envelope, err := simplecrypto.Seal(innerMsgBytes, &m.ownProfile.PrivateKey, &target.PublicKey)
	if err != nil {
		return middleware_entities.PendingMessage{}, fmt.Errorf("failed to seal message: %w", err)
	}
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return middleware_entities.PendingMessage{}, fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// 3. Create the P2P wire message (PendingMessage).
	msgTarget := middleware_entities.NewMessageTarget(message.ID, target.ID)
	pendingMsg := middleware_entities.NewPendingMessage(msgTarget, envelopeBytes)

	return pendingMsg, nil
}

// storeAsPending saves a message to the cache.
func (m *Middleware) storeAsPending(pendingMsg middleware_entities.PendingMessage) error {
	if err := m.messageCache.StoreImmune(pendingMsg); err != nil {
		return fmt.Errorf("failed to store pending message: %w", err)
	}
	return nil
}

func (m *Middleware) getUserIfConnected(userID entities.UserID) (net.IP, error) {
	data, err := m.getUserData(ezlog.NewLoggerContext("stump"), []entities.UserID{userID}, true)
	if err != nil {
		return nil, err
	}

	if len(data) != 1 {
		return nil, fmt.Errorf("user not found")
	}

	return *data[userID].IP, nil
}
