package idmanagerregisterer

import (
	"bytes"
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezcrypto/rsa"
	"eaglechat/common/ezlog"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type requestBody struct {
	Username  string `json:"username"`
	PublicKey []byte `json:"public_key"`
}

type responseBody struct {
	ID string `json:"id"`
}

// performHTTPRequest sends the final HTTP registration request.
func (r *registererImpl) performHTTPRequest(ctx context.Context, username string, pk *rsa.PublicKey, idManagerIP string) (entities.User, error) {
	ezlog.Log(ctx).Infof("Performing HTTP registration request to ID Manager at %s:%s", idManagerIP, r.idManagerPort)

	url := fmt.Sprintf("http://%s:%s/users/register", idManagerIP, r.idManagerPort)

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

	client := &http.Client{}
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

	return entities.NewUser(res.ID, username, *pk), nil
}
