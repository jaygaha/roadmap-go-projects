package backup

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

func mongoURI(dbConfig config.DatabaseConfig) string {
	authDB := dbConfig.AuthDB
	if authDB == "" {
		authDB = "admin"
	}

	if dbConfig.User != "" && dbConfig.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, authDB)
	}
	return fmt.Sprintf("mongodb://%s:%d/%s", dbConfig.Host, dbConfig.Port, dbConfig.Name)
}

// BackupMongoDB performs a MongoDB backup using mongodump
func BackupMongoDB(dbConfig config.DatabaseConfig, filePath string) error {
	zap.L().Info("Starting MongoDB backup",
		zap.String("database", dbConfig.Name),
		zap.String("host", dbConfig.Host),
		zap.Int("port", dbConfig.Port),
	)

	args := []string{
		fmt.Sprintf("--uri=%s", mongoURI(dbConfig)),
		fmt.Sprintf("--archive=%s", filePath),
		"--gzip",
	}

	cmd := exec.Command("mongodump", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongodump failed: %w\nstderr: %s", err, stderr.String())
	}

	zap.L().Info("MongoDB backup completed",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", filePath),
	)
	return nil
}
