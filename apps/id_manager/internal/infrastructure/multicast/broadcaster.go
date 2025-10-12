package multicast

import (
	"fmt"
	"os"
	"strconv"

	"eaglechat/apps/id_manager/internal/domain/entities"
	multicast "eaglechat/common/multicast/interface"

	"github.com/joho/godotenv"
)

type Broadcaster struct {
	multicastNet multicast.MulticastNetwork
	ip           string
}

func NewBroadcaster(multicastNet multicast.MulticastNetwork) (*Broadcaster, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	ip := os.Getenv("ID_MANAGER_IP")
	if ip == "" {
		return nil, fmt.Errorf("ID_MANAGER_IP not found in .env file")
	}
	return &Broadcaster{multicastNet: multicastNet, ip: ip}, nil
}

func (b *Broadcaster) BroadcastUserUpdate(user *entities.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	msg, err := multicast.BuildIDManagerUpdate(
		user.ID,
		b.ip,
		strconv.Itoa(8080), // manager's port
		user.PublicKeyPEM,
	)
	if err != nil {
		return fmt.Errorf("failed to build update message: %w", err)
	}

	return b.multicastNet.Broadcast(msg)
}
