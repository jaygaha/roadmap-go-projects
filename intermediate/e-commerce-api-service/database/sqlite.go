package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// New opens (or creates) the SQLite database and applies pragmas
// for better concurrent-read performance and data safety.
func New(dbPath string) (*sql.DB, error) {
	log.Printf("data/" + dbPath + "?_journal_mode=WAL&_foreign_keys=ON")
	db, err := sql.Open("sqlite3", "data/"+dbPath+"?_journal_mode=WAL&_foreign_keys=ON")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Verify connection is open.
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	log.Printf("[DB] Successfully connected to SQLite database: %s", dbPath)

	return db, nil
}
