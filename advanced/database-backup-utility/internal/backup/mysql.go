package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/utils"
	"go.uber.org/zap"
)

// BackupMySQL performs a MySQL backup using mysqldump
func BackupMySQL(dbConfig config.DatabaseConfig, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer file.Close()

	cmd := exec.Command("mysqldump",
		"-h", dbConfig.Host,
		"-P", fmt.Sprintf("%d", dbConfig.Port),
		"-u", dbConfig.User,
		dbConfig.Name,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", dbConfig.Password))
	cmd.Stdout = file

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return utils.HandleError(fmt.Errorf("mysqldump failed: %v", err), "MySQL backup command failed",
			zap.String("database", dbConfig.Name),
			zap.String("host", dbConfig.Host),
			zap.Int("port", dbConfig.Port),
			zap.String("stderr", stderr.String()),
		)
	}

	zap.L().Info("MySQL backup completed successfully",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", filePath),
	)
	return nil
}
