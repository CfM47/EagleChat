package idmanagerconn

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

type getRandomUsersResponse struct {
	Users []userData `json:"users"`
}

func (c *idManagerConnectionImpl) GetRandomConnectedUsers(ctx context.Context, count int) ([]middleware_entities.UserData, error) {
	url := fmt.Sprintf("%s/users/random?amount=%d", c.baseURL, count)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to create request to get random users: %v", err)
		return nil, err
	}
	req.Header.Add("X-Client-ID", string(c.ownID))

	resp, err := c.client.Do(req)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to get random users: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ezlog.Log(ctx).Errorf("Failed to get random users: received status %s", resp.Status)
		return nil, fmt.Errorf("failed to get random users: %s", resp.Status)
	}

	var responseData getRandomUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		ezlog.Log(ctx).Errorf("Failed to decode get random users response: %v", err)
		return nil, err
	}

	var result []middleware_entities.UserData
	for _, u := range responseData.Users {
		publicKey, err := rsa.PublicKeyFromBytes(u.PublicKey)
		if err != nil {
			ezlog.Log(ctx).Warnf("Invalid public key found while getting random user '%s' from ID manager", u.ID)
			continue
		}

		ip := net.ParseIP(u.IP)
		if ip == nil {
			ezlog.Log(ctx).Warnf("Invalid IP address found while getting random user '%s' from ID manager: %s", u.ID, u.IP)
			continue
		}

		userData := middleware_entities.NewUserData(
			entities.NewUser(
				u.ID,
				u.Username,
				*publicKey,
				u.LastSeen,
			),
			&ip,
		)

		result = append(result, userData)
	}

	return result, nil
}
