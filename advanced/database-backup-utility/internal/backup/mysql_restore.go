package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// RestoreMySQL performs a MySQL restore using mysql command
func RestoreMySQL(dbConfig config.DatabaseConfig, backupFilePath string) error {
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}

	backupData, err := os.ReadFile(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %v", err)
	}

	cmd := exec.Command("mysql",
		"-h", dbConfig.Host,
		"-P", fmt.Sprintf("%d", dbConfig.Port),
		"-u", dbConfig.User,
		dbConfig.Name,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", dbConfig.Password))
	cmd.Stdin = bytes.NewReader(backupData)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore failed: %v\nError: %s", err, stderr.String())
	}

	zap.L().Sugar().Infof("MySQL restore completed for %s from %s", dbConfig.Name, backupFilePath)
	return nil
}
