package backup

import (
	"fmt"
	"os"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// RestoreSQLite restores a SQLite database by replacing the database file with the backup
func RestoreSQLite(dbConfig config.DatabaseConfig, backupFilePath string) error {
	dbPath := dbConfig.Host
	if dbPath == "" {
		return fmt.Errorf("SQLite database path is empty; set 'host' to the .db file path in config")
	}

	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}

	zap.L().Info("Starting SQLite restore",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", backupFilePath),
		zap.String("target", dbPath),
	)

	if err := copyFile(backupFilePath, dbPath); err != nil {
		return fmt.Errorf("SQLite restore (file copy) failed: %w", err)
	}

	zap.L().Info("SQLite restore completed",
		zap.String("database", dbConfig.Name),
		zap.String("restored_to", dbPath),
	)
	return nil
}
