package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// RestoreMongoDB restores a full MongoDB database from a mongodump archive
func RestoreMongoDB(dbConfig config.DatabaseConfig, backupFilePath string) error {
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}

	zap.L().Info("Starting MongoDB full restore",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", backupFilePath),
	)

	args := []string{
		fmt.Sprintf("--uri=%s", mongoURI(dbConfig)),
		fmt.Sprintf("--archive=%s", backupFilePath),
		"--gzip",
		"--drop",
	}

	cmd := exec.Command("mongorestore", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongorestore failed: %w\nstderr: %s", err, stderr.String())
	}

	zap.L().Info("MongoDB full restore completed",
		zap.String("database", dbConfig.Name),
	)
	return nil
}

// RestoreMongoDBSelective restores a single collection from a mongodump archive
func RestoreMongoDBSelective(dbConfig config.DatabaseConfig, backupFilePath string, collection string) error {
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}

	zap.L().Info("Starting MongoDB selective restore",
		zap.String("database", dbConfig.Name),
		zap.String("collection", collection),
		zap.String("backup_file", backupFilePath),
	)

	args := []string{
		fmt.Sprintf("--uri=%s", mongoURI(dbConfig)),
		fmt.Sprintf("--archive=%s", backupFilePath),
		"--gzip",
		fmt.Sprintf("--nsInclude=%s.%s", dbConfig.Name, collection),
		"--drop",
	}

	cmd := exec.Command("mongorestore", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongorestore selective failed: %w\nstderr: %s", err, stderr.String())
	}

	zap.L().Info("MongoDB selective restore completed",
		zap.String("database", dbConfig.Name),
		zap.String("collection", collection),
	)
	return nil
}
