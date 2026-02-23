package database

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestInitDBAndSeeds(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "tracker_test.db")
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("query error: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected seed exercises inserted")
	}
}
