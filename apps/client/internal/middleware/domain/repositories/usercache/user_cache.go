package usercache

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"errors"
)

type UserCacheRepository interface {
	Save(middleware_entities.UserData) error
	Get(entities.UserID) (middleware_entities.UserData, error)
}

var ErrUserNotFound = errors.New("user not found")
