package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the shared database instance.
var DB *sql.DB

const dbPath = "./todo.db"

// InitDB initializes the database connection.
// It creates a database if it doesn't exist.
func InitDB() {
	var err error
	// Check if the database file exists
	// If not, create it
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Could not open database: %s\n", err.Error())
	}

	// Test the database connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Could not connect to database: %s\n", err.Error())
	}

	// Run migrations
	runMigrations()

	log.Println("Database connection established")
}

// runMigrations runs the database migrations.
func runMigrations() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	todoTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		is_completed BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	if _, err := DB.Exec(userTable); err != nil {
		log.Fatalf("Could not create users table: %s\n", err.Error())
	}
	if _, err := DB.Exec(todoTable); err != nil {
		log.Fatalf("Could not create todos table: %s\n", err.Error())
	}

	log.Println("Database migrations completed")
}
