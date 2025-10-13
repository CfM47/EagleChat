package usercache

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"errors"
	"time"
)

const DefaultIPExpirationTime time.Duration = time.Second * 30

// UserCacheRepository stores user data for each user that the client knows about
// it is expected to be thread-safe
type UserCacheRepository interface {
	// Save will store a user's data into the cache, updating the user, and its IP
	// expiration timer, when a user with the same id but different public key to
	// some other is inserted, an error should be logged, but not returned
	Save(middleware_entities.UserData) error

	// Get will return a user's data, nullifying the IP after a given time from the
	// user's last seen field
	Get(entities.UserID) (middleware_entities.UserData, error)
}

var ErrUserNotFound = errors.New("user not found")
