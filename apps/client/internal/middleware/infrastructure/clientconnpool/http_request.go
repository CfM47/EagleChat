package clientconnpool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto"
	"eaglechat/common/simplecrypto/rsa"

	"github.com/google/uuid"
)

func (c *clientConnPoolImpl) sendMessageRequest(ctx context.Context, messages []entities.PendingMessage, targetIP net.IP, targetPK rsa.PublicKey) error {
	ezlog.Log(ctx).Infof("Sending message request with %d messages to client at %s", len(messages), targetIP)

	nonce := createNonce()

	request := newMessageRequest(messages, nonce)
	encodedRequest, err := json.Marshal(request)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to marshall message request: %v", err)
		return err
	}

	envelope, err := simplecrypto.Seal(encodedRequest, &c.ownProfile.PrivateKey, &targetPK)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to create secure envelope: %v", err)
		return err
	}

	responseEnvelope, err := c.performHTTPRequest(ctx, targetIP, *envelope)
	if err != nil {
		ezlog.Log(ctx).Warnf("Failed to perform http request: %v", err)
		return err
	}

	nonceBytes, _, err := simplecrypto.Open(responseEnvelope, &c.ownProfile.PrivateKey)
	if err != nil {
		ezlog.Log(ctx).Warnf("Message target failed challenge: %v", err)
		return services.ErrAuthenticationFailed
	}

	if string(nonceBytes) != nonce {
		ezlog.Log(ctx).Warnf("Message target sent wrong nonce (expected: '%s' | received: '%s')", nonce, string(nonceBytes))
		return services.ErrAuthenticationFailed
	}

	return nil
}

func (c *clientConnPoolImpl) performHTTPRequest(ctx context.Context, targetIP net.IP, requestBody simplecrypto.SecureEnvelope) (*simplecrypto.SecureEnvelope, error) {
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("http://%s:%d/message", targetIP.String(), c.listenPort)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseEnvelope simplecrypto.SecureEnvelope
	if err := json.Unmarshal(responseBody, &responseEnvelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response envelope: %w", err)
	}

	return &responseEnvelope, nil
}

func createNonce() string {
	return uuid.NewString()
}
