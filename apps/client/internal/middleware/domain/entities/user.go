package entities

import (
	"net"

	"eaglechat/apps/client/internal/domain/entities"
)

type UserData struct {
	entities.User

	IP *net.IP
}

func NewUserData(user entities.User, IP *net.IP) UserData {
	return UserData{
		User: user,
		IP:   IP,
	}
}

func (d *UserData) GetUser() entities.User {
	return entities.NewUser(string(d.ID), d.Name, d.PublicKey)
}
