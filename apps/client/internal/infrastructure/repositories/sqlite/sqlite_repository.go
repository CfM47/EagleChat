package sqliterepository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"eaglechat/apps/client/internal/domain/repositories"

	_ "modernc.org/sqlite"
)

// sqliteRepository implements the unified ClientRepository interface.
type sqliteRepository struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteRepository creates a new SQLite-based repository and initializes the schema.
func NewSQLiteRepository(path string) (repositories.ClientRepository, error) {
	// Ensure the directory for the database file exists.
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	repo := &sqliteRepository{db: db}

	if err := repo.initSchema(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *sqliteRepository) initSchema() error {
	schema := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		public_key BLOB NOT NULL,
		last_seen INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS own_profile (
		user_id TEXT PRIMARY KEY,
		private_key BLOB NOT NULL,
		FOREIGN KEY(user_id) REFERENCES user(id)
	);

	CREATE TABLE IF NOT EXISTS chat (
		partner_id TEXT PRIMARY KEY,
		unread_count INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY(partner_id) REFERENCES user(id)
	);

	CREATE TABLE IF NOT EXISTS message (
		message_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		partner_id TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		PRIMARY KEY (message_id, sender_id),
		FOREIGN KEY(partner_id) REFERENCES chat(partner_id)
	);

	CREATE INDEX IF NOT EXISTS idx_message_chat_time ON message (partner_id, timestamp);
	`

	_, err := r.db.Exec(schema)
	return err
}
