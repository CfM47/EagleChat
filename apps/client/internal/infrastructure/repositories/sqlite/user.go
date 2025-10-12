package sqliterepository

import (
	"database/sql"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"
	"eaglechat/common/simplecrypto/rsa"
	"time"
)

// SaveUser saves public data about a user.
func (r *sqliteRepository) SaveUser(user entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.saveUser(tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

// saveUser is the internal, non-locking implementation for saving a user's public data.
// It requires a transaction to be passed in so it can be part of a larger atomic operation.
func (r *sqliteRepository) saveUser(tx *sql.Tx, user entities.User) error {
	pubKeyBytes, err := user.PublicKey.ToBytes()
	if err != nil {
		return err
	}

	query := `INSERT OR REPLACE INTO user (id, name, public_key, last_seen) VALUES (?, ?, ?, ?);`
	_, err = tx.Exec(query, user.ID, user.Name, pubKeyBytes, user.LastSeen.Unix())
	return err
}

// GetUser retrieves public data about a user.
func (r *sqliteRepository) GetUser(userID entities.UserID) (entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := `SELECT id, name, public_key, last_seen FROM user WHERE id = ?;`
	row := r.db.QueryRow(query, userID)

	var user entities.User
	var pubKeyBytes []byte
	var lastSeenUnix int64

	err := row.Scan(&user.ID, &user.Name, &pubKeyBytes, &lastSeenUnix)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, repositories.ErrUserNotFound
		}
		return entities.User{}, err
	}

	pubKey, err := rsa.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return entities.User{}, err
	}
	user.PublicKey = *pubKey
	user.LastSeen = time.Unix(lastSeenUnix, 0)

	return user, nil
}
