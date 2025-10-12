package sqliterepository

import (
	"database/sql"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

// SaveOwnProfile saves the full user profile, including the private key, to local storage.
func (r *sqliteRepository) SaveOwnProfile(profile entities.OwnProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback on error

	// Save the public part of the user
	user := profile.User
	pubKeyBytes, err := user.PublicKey.ToBytes()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT OR REPLACE INTO user (id, name, public_key, last_seen) VALUES (?, ?, ?, ?);`,
		user.ID, user.Name, pubKeyBytes, user.LastSeen.Unix(),
	)
	if err != nil {
		return err
	}

	// Save the private key
	_, err = tx.Exec(
		`INSERT OR REPLACE INTO own_profile (user_id, private_key) VALUES (?, ?);`,
		user.ID, profile.PrivateKey.ToBytes(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetOwnProfile retrieves the full user profile from local storage.
func (r *sqliteRepository) GetOwnProfile() (entities.OwnProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getOwnProfile()
}

// getOwnProfile is the internal, non-locking implementation for retrieving the user profile.
// It is called by public methods that have already acquired the appropriate lock.
func (r *sqliteRepository) getOwnProfile() (entities.OwnProfile, error) {
	query := `SELECT u.id, u.name, u.public_key, u.last_seen, p.private_key
			 FROM own_profile p
			 JOIN user u ON p.user_id = u.id;`

	row := r.db.QueryRow(query)

	var profile entities.OwnProfile
	var pubKeyBytes, privKeyBytes []byte
	var lastSeenUnix int64

	err := row.Scan(&profile.User.ID, &profile.User.Name, &pubKeyBytes, &lastSeenUnix, &privKeyBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.OwnProfile{}, repositories.ErrOwnProfileNotSet
		}
		return entities.OwnProfile{}, err
	}

	// Deserialize keys
	pubKey, err := rsa.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return entities.OwnProfile{}, err
	}
	profile.User.PublicKey = *pubKey

	privKey, err := rsa.PrivateKeyFromBytes(privKeyBytes)
	if err != nil {
		return entities.OwnProfile{}, err
	}
	profile.PrivateKey = *privKey

	return profile, nil
}
