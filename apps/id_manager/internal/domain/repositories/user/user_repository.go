package user

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"errors"
	"net"
)

type UserRepository interface {
	Save(user *entities.User) error
	FindByID(ID string) (*entities.User, error)
	FindAll() ([]*entities.User, error)
	Delete(ID string) error
	UpdateIP(ID string, ip net.IP) error
	Update(user *entities.User) error
	Create(user *entities.User) (*entities.User, error)
}

var ErrUserNotFound = errors.New("user not found")
