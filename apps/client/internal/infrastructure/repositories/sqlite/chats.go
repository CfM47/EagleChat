package sqliterepository

import (
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

// GetAllChatOverviews returns the corresponding ChatOverview object for each chat.
func (r *sqliteRepository) GetAllChatOverviews() ([]repositories.ChatOverview, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := `
		SELECT
			c.partner_id,
			c.unread_count,
			latest_message.content,
			latest_message.timestamp,
			u.name,
			u.public_key,
			u.last_seen
		FROM chat c
		JOIN user u ON c.partner_id = u.id
		JOIN (
			SELECT 
				partner_id,
				content,
				timestamp
			FROM message
			WHERE (partner_id, timestamp) IN (
				SELECT partner_id, MAX(timestamp) FROM message GROUP BY partner_id
			)
		) AS latest_message ON c.partner_id = latest_message.partner_id
		ORDER BY latest_message.timestamp DESC;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overviews []repositories.ChatOverview
	for rows.Next() {
		var overview repositories.ChatOverview
		var timestampUnix, lastSeenUnix int64
		var pubKeyBytes []byte

		err := rows.Scan(
			&overview.Partner.ID, &overview.UnreadCount, &overview.LastMessage, &timestampUnix,
			&overview.Partner.Name, &pubKeyBytes, &lastSeenUnix,
		)
		if err != nil {
			return nil, err
		}

		overview.Timestamp = time.Unix(timestampUnix, 0)
		overview.Partner.LastSeen = time.Unix(lastSeenUnix, 0)

		pubKey, err := rsa.PublicKeyFromBytes(pubKeyBytes)
		if err != nil {
			return nil, err
		}
		overview.Partner.PublicKey = *pubKey

		overviews = append(overviews, overview)
	}

	return overviews, rows.Err()
}

// IncrementUnreadCount increments the unread message count for a given chat.
func (r *sqliteRepository) IncrementUnreadCount(partnerID entities.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `UPDATE chat SET unread_count = unread_count + 1 WHERE partner_id = ?;`
	_, err := r.db.Exec(query, partnerID)
	return err
}

// ResetUnreadCount resets the unread message count for a given chat to zero.
func (r *sqliteRepository) ResetUnreadCount(partnerID entities.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `UPDATE chat SET unread_count = 0 WHERE partner_id = ?;`
	_, err := r.db.Exec(query, partnerID)
	return err
}
