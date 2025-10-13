package entities

import (
	"eaglechat/apps/client/internal/domain/entities"
	"net"
)

type UserData struct {
	entities.User

	IP *net.IP
}

func (d *UserData) GetUser() entities.User {
	return entities.NewUser(string(d.ID), d.Name, d.PublicKey)
}
