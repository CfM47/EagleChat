package repositories

import (
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
)

type (
	UserRepository = user.UserRepository
)

var (
	ErrUserNotFound = user.ErrUserNotFound
)
