package sqliterepository

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/simplecrypto/rsa"
	"time"
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

	var partnerID entities.UserID
	if message.Sender.ID == ownProfile.User.ID {
		partnerID = message.Target.ID
	} else {
		partnerID = message.Sender.ID
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure the sender and target users exist in the user table
	if err := r.saveUser(tx, message.Sender); err != nil {
		return err
	}
	if err := r.saveUser(tx, message.Target); err != nil {
		return err
	}

	// Ensure the chat metadata row exists
	_, err = tx.Exec(`INSERT OR IGNORE INTO chat (partner_id) VALUES (?);`, partnerID)
	if err != nil {
		return err
	}

	// Insert the message, ignoring if it's a duplicate
	_, err = tx.Exec(
		`INSERT OR IGNORE INTO message (message_id, sender_id, partner_id, content, timestamp) VALUES (?, ?, ?, ?, ?);`,
		message.ID, message.Sender.ID, partnerID, message.Content, message.CreatedTime.Unix(),
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
