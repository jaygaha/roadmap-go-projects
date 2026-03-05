package backup

import (
	"fmt"
	"io"
	"os"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// BackupSQLite performs a SQLite backup by copying the database file
func BackupSQLite(dbConfig config.DatabaseConfig, filePath string) error {
	dbPath := dbConfig.Host
	if dbPath == "" {
		return fmt.Errorf("SQLite database path is empty; set 'host' to the .db file path in config")
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("SQLite database file not found: %s", dbPath)
	}

	zap.L().Info("Starting SQLite backup",
		zap.String("database", dbConfig.Name),
		zap.String("source", dbPath),
		zap.String("destination", filePath),
	)

	if err := copyFile(dbPath, filePath); err != nil {
		return fmt.Errorf("SQLite backup (file copy) failed: %w", err)
	}

	zap.L().Info("SQLite backup completed",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", filePath),
	)
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	tmpPath := dst + ".tmp"
	dstFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file %s: %w", tmpPath, err)
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		dstFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("file copy failed: %w", err)
	}

	if err := dstFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to flush destination file: %w", err)
	}

	if err := os.Rename(tmpPath, dst); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to finalise backup file: %w", err)
	}

	return nil
}
