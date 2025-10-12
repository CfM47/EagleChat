package middleware

import (
	"bytes"
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto"
	"encoding/json"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	messagecache "eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
)

// routeIncomingMessages is a background goroutine that processes all messages
// received from the P2P connection pool.
func (m *Middleware) routeIncomingMessages(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting incoming message router...")
	defer ezlog.Log(ctx).Info("Stopped incoming message router.")

	for rawMsg := range m.p2pConnPool.Receive() {
		var pendingMsg middleware_entities.PendingMessage
		if err := json.Unmarshal(rawMsg, &pendingMsg); err != nil {
			ezlog.Log(ctx).Errorf("Failed to unmarshal incoming message: %v", err)
			continue
		}

		if pendingMsg.Target.TargetID == m.ownUser.ID {
			ezlog.Log(ctx).Infof("Received message for self: %s", pendingMsg.Target.MessageID)

			messageCtx := ezlog.NewLoggerContext("message-for-self-hadler")
			m.handleMessageForSelf(messageCtx, pendingMsg)
		} else {
			ezlog.Log(ctx).Infof("Received message for other user %s: %s", pendingMsg.Target.TargetID, pendingMsg.Target.MessageID)

			messageCtx := ezlog.NewLoggerContext("message-for-other-handler")
			m.handleMessageForOther(messageCtx, pendingMsg)
		}
	}
}

// handleMessageForSelf processes a message that is intended for the current user.
// It decrypts, verifies, and forwards the message to the application.
func (m *Middleware) handleMessageForSelf(ctx context.Context, pendingMsg middleware_entities.PendingMessage) {
	var envelope simplecrypto.SecureEnvelope
	if err := json.Unmarshal(pendingMsg.Content, &envelope); err != nil {
		ezlog.Log(ctx).Errorf("Failed to unmarshal secure envelope for own message: %v", err)
		return
	}

	plaintext, senderPubKey, err := simplecrypto.Open(&envelope, &m.sk)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to open secure envelope for own message: %v", err)
		return
	}

	var msg entities.Message
	if err := json.Unmarshal(plaintext, &msg); err != nil {
		ezlog.Log(ctx).Errorf("Failed to unmarshal inner message for own message: %v", err)
		return
	}

	senderPubKeyBytes, err := senderPubKey.ToBytes()
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to serialize sender public key for own message: %v", err)
		return
	}
	msgSenderPubKeyBytes, err := msg.Sender.PublicKey.ToBytes()
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to serialize message sender public key for own message: %v", err)
		return
	}
	if !bytes.Equal(senderPubKeyBytes, msgSenderPubKeyBytes) {
		ezlog.Log(ctx).Errorf("Security alert: sender public key mismatch in message %s", pendingMsg.Target.MessageID)
		return
	}

	m.receivedMessages <- msg
}

// handleMessageForOther processes a message that is intended for another user.
// It stores the message in the cache for later forwarding.
func (m *Middleware) handleMessageForOther(ctx context.Context, pendingMsg middleware_entities.PendingMessage) {
	if err := m.messageCache.StoreExpiring(pendingMsg, messagecache.DefaultImmunityPeriod); err != nil {
		ezlog.Log(ctx).Errorf("Failed to store message for other user %s: %v", pendingMsg.Target.TargetID, err)
	}
	// TODO: Notify ID Manager that we have a pending message for another user.
}
