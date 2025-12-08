package clientconnpool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto"
)

func (c *clientConnPoolImpl) serve() {
	serveCtx := ezlog.WithNewLogger(context.Background(), "client-conn-pool-server")
	defer close(c.done)

	mux := http.NewServeMux()
	mux.HandleFunc("/message", c.HandleMessage)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.listenPort),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			ezlog.Log(serveCtx).Errorf("HTTP server failed: %v", err)
		}
	}()

	<-c.quit // Wait for shutdown signal

	shutdownCtx := ezlog.WithNewLogger(serveCtx, "http-server-shutdown")
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, shutdownTimeout) // Use the constant here
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		ezlog.Log(shutdownCtx).Warnf("HTTP server shutdown error: %v", err)
	}

	close(c.messages)
}

func (c *clientConnPoolImpl) HandleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := ezlog.WithNewLogger(r.Context(), "receive-message")
	ezlog.Log(ctx).Info("Received message request")

	envelope, err := unmarshallSecureEnvelope(r)
	if err != nil {
		msg := fmt.Sprintf("Invalid request body: %v", err)
		ezlog.Log(ctx).Warn(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	messageRequestBytes, senderPK, err := simplecrypto.Open(envelope, &c.ownProfile.PrivateKey)
	if err != nil {
		msg := fmt.Sprintf("Invalid secure envelope: %v", err)
		ezlog.Log(ctx).Warn(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	messageRequest, err := unmarshallRequestBytes(messageRequestBytes)
	if err != nil {
		msg := fmt.Sprintf("Invalid secure envelope content: %v", err)
		ezlog.Log(ctx).Warn(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	c.addMessagesToChannel(ctx, messageRequest.Messages)

	nonce := messageRequest.Nonce
	nonceEnvelope, err := simplecrypto.Seal([]byte(nonce), &c.ownProfile.PrivateKey, senderPK)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to create challenge response envelope: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	respond(ctx, nonceEnvelope, w)
}

func (c *clientConnPoolImpl) addMessagesToChannel(ctx context.Context, messages []entities.PendingMessage) {
	select {
	case c.messages <- messages:
		// Message sent immediately.
		return
	case <-c.quit:
		ezlog.Log(ctx).Warn("Pool is shutting down, dropping incoming message.")
		return
	default:
		// Channel is busy, spawn a shutdown-aware goroutine.
		go func() {
			select {
			case c.messages <- messages:
				// Message sent after a wait.
			case <-c.quit:
				ezlog.Log(ctx).Warn("Pool is shutting down, dropping incoming message.")
			}
		}()
	}
}

func unmarshallSecureEnvelope(r *http.Request) (*simplecrypto.SecureEnvelope, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var envelope simplecrypto.SecureEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secure envelope: %w", err)
	}

	return &envelope, nil
}

func unmarshallRequestBytes(messageRequestBytes []byte) (messageRequest, error) {
	var request messageRequest
	if err := json.Unmarshal(messageRequestBytes, &request); err != nil {
		return messageRequest{}, fmt.Errorf("failed to unmarshal message request: %w", err)
	}
	return request, nil
}

func respond(ctx context.Context, nonceEnvelope *simplecrypto.SecureEnvelope, w http.ResponseWriter) {
	responseBytes, err := json.Marshal(nonceEnvelope)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to marshal nonce envelope: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseBytes); err != nil {
		ezlog.Log(ctx).Warnf("Failed to write response: %v", err)
	}
}
