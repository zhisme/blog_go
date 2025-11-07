package repositories

import (
	"backend-go/internal/dto"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteMailingListRepository struct {
	db *sql.DB
}

func NewSqliteMailingListRepository(dbPath string) (*SqliteMailingListRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	repo := &SqliteMailingListRepository{db: db}

	// Initialize schema
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *SqliteMailingListRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS mailing_list (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_mailing_list_email ON mailing_list(email);
	CREATE INDEX IF NOT EXISTS idx_mailing_list_created_at ON mailing_list(created_at);
	`

	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func (r *SqliteMailingListRepository) Save(mailingList *dto.MailingList) error {
	createdAt := mailingList.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	query := `INSERT INTO mailing_list (username, email, created_at) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, mailingList.Username, mailingList.Email, createdAt)
	if err != nil {
		// Check if it's a unique constraint violation
		if err.Error() == "UNIQUE constraint failed: mailing_list.email" {
			log.Printf("Email already subscribed: %s", mailingList.Email)
			return nil // Same behavior as CSV implementation
		}
		return fmt.Errorf("failed to save mailing list entry: %w", err)
	}

	return nil
}

func (r *SqliteMailingListRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
