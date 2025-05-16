package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// variable to store database connection
var DB *sql.DB

// ConnectDB connects to the database
func ConnectDB() {
	var err error
	DB, err = sql.Open("sqlite", "./database/blog.db")

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// check if database is connected
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database", err)
	}

	log.Println("Connected to database")

	// Run migrations
	runMigrations()
}

// runMigrations runs the database migrations if needed
func runMigrations() {
	sql := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT, -- Auto-increment primary key
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			category TEXT NOT NULL,
			tags TEXT NOT NULL,
			created_at DATETIME NOT NULL,
        	updated_at DATETIME NULL
		);
	`
	_, err := DB.Exec(sql)
	if err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	log.Println("Migrations ran successfully")
}
