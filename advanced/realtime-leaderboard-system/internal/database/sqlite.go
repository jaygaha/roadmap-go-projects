package database

import (
	"log"

	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitSQLite() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.SQLitePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}
	log.Println("Connected to SQLite")
}
