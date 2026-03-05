package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// RestorePostgreSQL performs a PostgreSQL restore using psql
func RestorePostgreSQL(dbConfig config.DatabaseConfig, backupFilePath string) error {
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}

	cmd := exec.Command("psql",
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"-f", backupFilePath,
		"--single-transaction",
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore failed: %v\nError: %s", err, stderr.String())
	}

	zap.L().Sugar().Infof("PostgreSQL restore completed for %s from %s", dbConfig.Name, backupFilePath)
	return nil
}
