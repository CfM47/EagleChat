package user

import "eaglechat/apps/client/internal/domain/entities"

type UserRepository interface {
	// Save saves public data about a user.
	Save(user entities.User) error
	// Get retrieves public data about a user.
	Get(userID entities.UserID) (entities.User, error)

	// SaveOwnProfile saves the full user profile, including the private key, to local storage.
	SaveOwnProfile(profile entities.OwnProfile) error
	// GetOwnProfile retrieves the full user profile from local storage.
	GetOwnProfile() (entities.OwnProfile, error)
}
