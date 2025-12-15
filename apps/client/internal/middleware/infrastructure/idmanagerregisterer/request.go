package idmanagerregisterer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

type requestBody struct {
	Username  string `json:"username"`
	PublicKey []byte `json:"public_key"`
}

type responseBody struct {
	ID string `json:"id"`
}

// requestRegistration sends the final HTTP registration request.
func (r *registererImpl) requestRegistration(ctx context.Context, username string, pk *rsa.PublicKey, idManagerIP string, client *http.Client) (entities.User, error) {
	if client == nil {
		client = &http.Client{}
	}

	url := fmt.Sprintf("http://%s:%d/users/register", idManagerIP, r.idManagerPort)

	ezlog.Log(ctx).Infof("Performing HTTP registration request to ID Manager at %s", url)

	pubKeyBytes, err := pk.ToBytes()
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to serialize public key: %w", err)
	}

	reqBody := requestBody{
		Username:  username,
		PublicKey: pubKeyBytes,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal request body: %v", err)
		ezlog.Log(ctx).Error(msg)
		return entities.User{}, errors.New(msg)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		msg := fmt.Sprintf("failed to create HTTP request: %v", err)
		ezlog.Log(ctx).Error(msg)
		return entities.User{}, errors.New(msg)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("HTTP request to %s failed: %v", url, err)
		ezlog.Log(ctx).Error(msg)
		return entities.User{}, errors.New(msg)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("ID Manager responded with status: %s", resp.Status)
		ezlog.Log(ctx).Error(msg)
		return entities.User{}, errors.New(msg)
	}

	var res responseBody
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		msg := fmt.Sprintf("failed to decode response body: %v", err)
		ezlog.Log(ctx).Error(msg)
		return entities.User{}, errors.New(msg)
	}

	ezlog.Log(ctx).Infof("Successfully registered with ID Manager. Received user ID: %s", res.ID)

	// FIXME: we shouldn't have LastSeen for the client itself
	return entities.NewUser(res.ID, username, *pk, time.Time{}), nil
}
