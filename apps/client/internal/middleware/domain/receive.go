package middleware

import (
	"bytes"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto"
	"encoding/json"
	"log"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	messagecache "eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
)

// routeIncomingMessages is a background goroutine that processes all messages
// received from the P2P connection pool.
func (m *Middleware) routeIncomingMessages() {
	log.SetPrefix("[RouteIncomingMessages] ")

	for rawMsg := range m.p2pConnPool.Receive() {
		var pendingMsg middleware_entities.PendingMessage
		if err := json.Unmarshal(rawMsg, &pendingMsg); err != nil {
			log.Printf("failed to unmarshal incoming message: %v", err)
			continue
		}

		if pendingMsg.Target.Target == m.ownUser.ID {
			m.handleMessageForSelf(pendingMsg)
		} else {
			m.handleMessageForOther(pendingMsg)
		}
	}
}

// handleMessageForSelf processes a message that is intended for the current user.
// It decrypts, verifies, and forwards the message to the application.
func (m *Middleware) handleMessageForSelf(pendingMsg middleware_entities.PendingMessage) {
	var envelope simplecrypto.SecureEnvelope
	if err := json.Unmarshal(pendingMsg.Content, &envelope); err != nil {
		log.Printf("failed to unmarshal secure envelope for own message: %v", err)
		return
	}

	plaintext, senderPubKey, err := simplecrypto.Open(&envelope, &m.sk)
	if err != nil {
		log.Printf("failed to open secure envelope: %v", err)
		return
	}

	var msg entities.Message
	if err := json.Unmarshal(plaintext, &msg); err != nil {
		log.Printf("failed to unmarshal inner message: %v", err)
		return
	}

	senderPubKeyBytes, err := senderPubKey.ToBytes()
	if err != nil {
		log.Printf("failed to serialize sender public key: %v", err)
		return
	}
	msgSenderPubKeyBytes, err := msg.Sender.PublicKey.ToBytes()
	if err != nil {
		log.Printf("failed to serialize message sender public key: %v", err)
		return
	}
	if !bytes.Equal(senderPubKeyBytes, msgSenderPubKeyBytes) {
		log.Printf("security alert: sender public key mismatch in message %s", pendingMsg.Target.ID)
		return
	}

	m.receivedMessages <- msg
}

// handleMessageForOther processes a message that is intended for another user.
// It stores the message in the cache for later forwarding.
func (m *Middleware) handleMessageForOther(pendingMsg middleware_entities.PendingMessage) {
	if err := m.messageCache.StoreExpiring(pendingMsg, messagecache.DefaultImmunityPeriod); err != nil {
		log.Printf("failed to store message for other user %s: %v", pendingMsg.Target.Target, err)
	}
	// TODO: Notify ID Manager that we have a pending message for another user.
}
