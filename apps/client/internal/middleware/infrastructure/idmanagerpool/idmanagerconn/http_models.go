package idmanagerconn

import (
	"context"
	"fmt"
	"net"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

type userData struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	PublicKey []byte    `json:"public_key"`
	IP        string    `json:"ip"`
	LastSeen  time.Time `json:"last_seen"`
}

func buildUserData(ctx context.Context, u userData) (middleware_entities.UserData, error) {
	publicKey, err := rsa.PublicKeyFromBytes(u.PublicKey)
	if err != nil {
		ezlog.Log(ctx).Warnf("Invalid public key found while querying user '%s' from ID manager", u.ID)
		return middleware_entities.UserData{}, nil
	}

	ip := net.ParseIP(u.IP)
	if ip == nil {
		ezlog.Log(ctx).Warnf("Invalid IP address found while querying user '%s' from ID manager: %s", u.ID, u.IP)
		return middleware_entities.UserData{}, fmt.Errorf("invalid IP address for user %s: %s", u.ID, u.IP)
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

	return userData, nil
}
