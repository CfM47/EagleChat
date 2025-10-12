package jsonusercache

import (
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
)

// Get will return a user's data, nullifying the IP after a given time from the
// user's last seen field
func (c *jsonUserCache) Get(userID entities.UserID) (middleware_entities.UserData, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cache, err := c.readCache()
	if err != nil {
		return middleware_entities.UserData{}, err
	}

	userData, ok := cache[string(userID)]
	if !ok {
		return middleware_entities.UserData{}, usercache.ErrUserNotFound
	}

	// Nullify the IP if it has expired.
	if time.Since(userData.LastSeen) > c.expirationTime {
		userData.IP = nil
	}

	return userData, nil
}
