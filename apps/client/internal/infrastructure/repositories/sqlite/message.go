package sqliterepository

import (
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/simplecrypto/rsa"
)

// SaveMessage stores a message in the repository, ignoring duplicates.
func (r *sqliteRepository) SaveMessage(message entities.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get own profile to determine the partner ID
	ownProfile, err := r.getOwnProfile()
	if err != nil {
		return err
	}

	var partner entities.User
	if message.Sender.ID == ownProfile.User.ID {
		partner = message.Target
	} else {
		partner = message.Sender
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure the partner user exist in the user table
	if err := r.saveUser(tx, partner); err != nil {
		return err
	}

	// Insert the message, ignoring if it's a duplicate
	_, err = tx.Exec(
		`INSERT OR IGNORE INTO message (message_id, sender_id, partner_id, content, timestamp) VALUES (?, ?, ?, ?, ?);`,
		message.ID, message.Sender.ID, partner.ID, message.Content, message.CreatedTime.Unix(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetChat returns all messages for a given chat, sorted by creation time.
func (r *sqliteRepository) GetChat(partnerID entities.UserID) ([]entities.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// First, check if the user exists. If not, we can't have a chat with them.
	_, err := r.GetUser(partnerID)
	if err != nil {
		return nil, err // This will be ErrUserNotFound if they don't exist
	}

	ownProfile, err := r.GetOwnProfile()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			m.message_id, m.content, m.timestamp,
			sender.id, sender.name, sender.public_key, sender.last_seen,
			target.id, target.name, target.public_key, target.last_seen
		FROM message m
		JOIN user sender ON m.sender_id = sender.id
		-- The target is the user in the chat who is NOT the sender
		JOIN user target ON (CASE WHEN m.sender_id = ? THEN ? ELSE m.sender_id END) = target.id
		WHERE m.partner_id = ?
		ORDER BY m.timestamp ASC;
	`

	rows, err := r.db.Query(query, ownProfile.User.ID, partnerID, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entities.Message
	for rows.Next() {
		var msg entities.Message
		var timestampUnix int64
		var senderPubKeyBytes, targetPubKeyBytes []byte
		var senderLastSeenUnix, targetLastSeenUnix int64

		err := rows.Scan(
			&msg.ID, &msg.Content, &timestampUnix,
			&msg.Sender.ID, &msg.Sender.Name, &senderPubKeyBytes, &senderLastSeenUnix,
			&msg.Target.ID, &msg.Target.Name, &targetPubKeyBytes, &targetLastSeenUnix,
		)
		if err != nil {
			return nil, err
		}

		msg.CreatedTime = time.Unix(timestampUnix, 0)
		msg.Sender.LastSeen = time.Unix(senderLastSeenUnix, 0)
		msg.Target.LastSeen = time.Unix(targetLastSeenUnix, 0)

		senderPubKey, err := rsa.PublicKeyFromBytes(senderPubKeyBytes)
		if err != nil {
			return nil, err
		}
		msg.Sender.PublicKey = *senderPubKey

		targetPubKey, err := rsa.PublicKeyFromBytes(targetPubKeyBytes)
		if err != nil {
			return nil, err
		}
		msg.Target.PublicKey = *targetPubKey

		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *sqliteRepository) MessageExists(messageID string, senderID entities.UserID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := `
		SELECT 1 FROM message WHERE message_id = ? AND sender_id = ? LIMIT 1;
	`

	rows, err := r.db.Query(query, messageID, senderID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
