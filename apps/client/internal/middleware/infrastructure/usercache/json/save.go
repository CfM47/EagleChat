package jsonusercache

import (
	"bytes"
	"log"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

// Save will store a user's data into the cache, updating the user, and its IP
// expiration timer, when a user with the same id but different public key to
// some other is inserted, an error should be logged, but not returned
func (c *jsonUserCache) Save(userData middleware_entities.UserData) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, err := c.readCache()
	if err != nil {
		return err
	}

	// Check for public key mismatch if the user already exists.
	if existingUser, ok := cache[string(userData.ID)]; ok {
		existingPK, _ := existingUser.PublicKey.ToBytes()
		newPK, _ := userData.PublicKey.ToBytes()
		if !bytes.Equal(existingPK, newPK) {
			log.Printf("Warning: saving user %s with a different public key.", userData.ID)
		}
	}

	// The LastSeen field is assumed to be set by the caller.
	cache[string(userData.ID)] = userData

	return c.writeCache(cache)
}
