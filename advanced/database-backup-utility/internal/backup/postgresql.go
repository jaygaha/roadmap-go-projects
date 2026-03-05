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

// BackupPostgreSQL performs a PostgreSQL backup using pg_dump
func BackupPostgreSQL(dbConfig config.DatabaseConfig, filePath string) error {
	cmd := exec.Command("pg_dump",
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"-f", filePath,
		"-w",
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return utils.HandleError(fmt.Errorf("pg_dump failed: %v", err), "PostgreSQL backup command failed",
			zap.String("database", dbConfig.Name),
			zap.String("host", dbConfig.Host),
			zap.Int("port", dbConfig.Port),
			zap.String("stderr", stderr.String()),
		)
	}

	zap.L().Info("PostgreSQL backup completed successfully",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", filePath),
	)
	return nil
}
