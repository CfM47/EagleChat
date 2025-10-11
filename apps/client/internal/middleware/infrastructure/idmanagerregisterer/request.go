package idmanagerregisterer

import (
	"bytes"
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"encoding/json"
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
		return entities.User{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return entities.User{}, fmt.Errorf("http request to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entities.User{}, fmt.Errorf("registration failed with status: %s", resp.Status)
	}

	var res responseBody
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return entities.User{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return entities.NewUser(res.ID, username, *pk), nil
}
