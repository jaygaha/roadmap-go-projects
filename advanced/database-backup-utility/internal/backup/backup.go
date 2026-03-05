package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/notification"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/storage"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/utils"
	"go.uber.org/zap"
)

// BackupDatabase performs a backup for the given database configuration
func BackupDatabase(dbConfig config.DatabaseConfig, backupConfig config.BackupConfig, storageConfig config.StorageConfig, notifConfig config.NotificationConfig) error {
	backupType := strings.ToLower(backupConfig.Type)
	if backupType == "" {
		backupType = "full"
	}

	ext := backupFileExtension(dbConfig.Type)
	filename := fmt.Sprintf("%s_%s_backup_%s%s", dbConfig.Name, backupType, time.Now().Format("20060102_150405"), ext)
	filePath := filepath.Join(storageConfig.Path, filename)

	if storageConfig.Type == "local" || storageConfig.Path != "" {
		if err := os.MkdirAll(storageConfig.Path, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory %s: %v", storageConfig.Path, err)
		}
	}

	zap.L().Info("Starting backup",
		zap.String("database", dbConfig.Name),
		zap.String("type", dbConfig.Type),
		zap.String("backup_type", backupType),
	)

	var state *BackupState
	if backupType != "full" {
		var err error
		state, err = LoadState(dbConfig.Name)
		if err != nil {
			return fmt.Errorf("failed to load backup state for %s: %v", dbConfig.Name, err)
		}
	}

	var newState *BackupState
	var backupErr error

	switch dbConfig.Type {
	case "postgres":
		newState, backupErr = runPostgreSQLBackup(dbConfig, filePath, backupType, state)
	case "mysql":
		newState, backupErr = runMySQLBackup(dbConfig, filePath, backupType, state)
	case "mongodb":
		newState, backupErr = runMongoDBBackup(dbConfig, filePath, backupType, state)
	case "sqlite":
		newState, backupErr = runSQLiteBackup(dbConfig, filePath, backupType, state)
	default:
		return utils.HandleError(
			fmt.Errorf("unsupported database type: %s", dbConfig.Type),
			"Database type not supported",
			zap.String("database", dbConfig.Name),
			zap.String("type", dbConfig.Type),
		)
	}

	if backupErr != nil {
		if notifConfig.Slack.Enabled && notifConfig.Slack.OnFailure {
			notification.NotifyFailure(notifConfig.Slack.WebhookURL, dbConfig.Name, "backup", backupErr)
		}
		return utils.HandleError(backupErr, "Database backup failed",
			zap.String("database", dbConfig.Name),
			zap.String("type", dbConfig.Type),
			zap.String("backup_file", filePath),
		)
	}

	if newState != nil {
		if backupType == "full" && newState == nil {
			newState = &BackupState{
				DBName:         dbConfig.Name,
				LastFullBackup: time.Now(),
				LastBackupTime: time.Now(),
				LastBackupType: "full",
			}
		}
		if err := SaveState(newState); err != nil {
			zap.L().Warn("Failed to save backup state",
				zap.String("database", dbConfig.Name),
				zap.Error(err),
			)
		}
	}

	finalFilePath := filePath
	if backupConfig.Compress {
		compressedFilePath := fmt.Sprintf("%s.gz", filePath)

		compressionLevel := backupConfig.CompressionLevel
		if compressionLevel == 0 {
			compressionLevel = 6
		}

		if err := compressBackupFile(filePath, compressedFilePath, compressionLevel); err != nil {
			return fmt.Errorf("failed to compress backup: %v", err)
		}

		finalFilePath = compressedFilePath
		defer os.Remove(filePath)
	}

	if storageConfig.Type == "local" && storageConfig.Retain > 0 {
		if err := enforceRetentionLocal(finalFilePath, storageConfig, true); err != nil {
			return fmt.Errorf("failed to enforce retention policy: %v", err)
		}
	}

	if err := uploadToStorage(finalFilePath, storageConfig, true); err != nil {
		return fmt.Errorf("failed to upload backup to storage: %v", err)
	}

	zap.L().Sugar().Infof("Backup for %s (%s) completed and uploaded to storage.", dbConfig.Name, backupType)
	if notifConfig.Slack.Enabled && notifConfig.Slack.OnSuccess {
		notification.NotifySuccess(notifConfig.Slack.WebhookURL, dbConfig.Name, backupType, finalFilePath)
	}
	return nil
}

func runPostgreSQLBackup(dbConfig config.DatabaseConfig, filePath, backupType string, state *BackupState) (*BackupState, error) {
	switch backupType {
	case "incremental":
		return BackupPostgreSQLIncremental(dbConfig, filePath, state)
	case "differential":
		return BackupPostgreSQLDifferential(dbConfig, filePath, state)
	default:
		if err := BackupPostgreSQL(dbConfig, filePath); err != nil {
			return nil, err
		}
		lsn, _ := getPGCurrentLSN(dbConfig)
		return &BackupState{
			DBName:         dbConfig.Name,
			LastFullBackup: time.Now(),
			LastBackupTime: time.Now(),
			LastBackupType: "full",
			PGLastLSN:      lsn,
		}, nil
	}
}

func runMySQLBackup(dbConfig config.DatabaseConfig, filePath, backupType string, state *BackupState) (*BackupState, error) {
	switch backupType {
	case "incremental":
		return BackupMySQLIncremental(dbConfig, filePath, state)
	case "differential":
		return BackupMySQLDifferential(dbConfig, filePath, state)
	default:
		if err := BackupMySQL(dbConfig, filePath); err != nil {
			return nil, err
		}
		file, pos, err := getMySQLBinlogPosition(dbConfig)
		newState := &BackupState{
			DBName:         dbConfig.Name,
			LastFullBackup: time.Now(),
			LastBackupTime: time.Now(),
			LastBackupType: "full",
		}
		if err == nil {
			newState.MySQLBinlogFile = file
			newState.MySQLBinlogPos = pos
		}
		return newState, nil
	}
}

func runMongoDBBackup(dbConfig config.DatabaseConfig, filePath, backupType string, state *BackupState) (*BackupState, error) {
	if backupType != "full" {
		zap.L().Warn("MongoDB incremental/differential not yet supported — running full backup",
			zap.String("database", dbConfig.Name),
		)
	}
	if err := BackupMongoDB(dbConfig, filePath); err != nil {
		return nil, err
	}
	return &BackupState{
		DBName:         dbConfig.Name,
		LastFullBackup: time.Now(),
		LastBackupTime: time.Now(),
		LastBackupType: "full",
	}, nil
}

func runSQLiteBackup(dbConfig config.DatabaseConfig, filePath, backupType string, state *BackupState) (*BackupState, error) {
	if backupType != "full" {
		zap.L().Warn("SQLite incremental/differential not supported — running full backup",
			zap.String("database", dbConfig.Name),
		)
	}
	if err := BackupSQLite(dbConfig, filePath); err != nil {
		return nil, err
	}
	return &BackupState{
		DBName:         dbConfig.Name,
		LastFullBackup: time.Now(),
		LastBackupTime: time.Now(),
		LastBackupType: "full",
	}, nil
}

func backupFileExtension(dbType string) string {
	switch dbType {
	case "mongodb":
		return ".archive"
	case "sqlite":
		return ".db"
	default:
		return ".sql"
	}
}

// RestoreDatabase restores a database from a backup file
func RestoreDatabase(dbConfig config.DatabaseConfig, backupFilePath string, tables []string, notifConfig config.NotificationConfig) error {
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}

	var err error
	switch dbConfig.Type {
	case "postgres":
		if len(tables) > 0 {
			err = RestorePostgreSQLSelective(dbConfig, backupFilePath, tables)
		} else {
			err = RestorePostgreSQL(dbConfig, backupFilePath)
		}
	case "mysql":
		if len(tables) > 0 {
			err = RestoreMySQLSelective(dbConfig, backupFilePath, tables)
		} else {
			err = RestoreMySQL(dbConfig, backupFilePath)
		}
	case "mongodb":
		if len(tables) > 0 {
			err = RestoreMongoDBSelective(dbConfig, backupFilePath, tables[0])
		} else {
			err = RestoreMongoDB(dbConfig, backupFilePath)
		}
	case "sqlite":
		err = RestoreSQLite(dbConfig, backupFilePath)
	default:
		return fmt.Errorf("unsupported database type for restore: %s", dbConfig.Type)
	}

	if err != nil {
		if notifConfig.Slack.Enabled && notifConfig.Slack.OnFailure {
			notification.NotifyFailure(notifConfig.Slack.WebhookURL, dbConfig.Name, "restore", err)
		}
		return fmt.Errorf("failed to restore database %s: %v", dbConfig.Name, err)
	}

	zap.L().Sugar().Infof("Restore for %s completed successfully", dbConfig.Name)
	if notifConfig.Slack.Enabled && notifConfig.Slack.OnSuccess {
		notification.NotifyRestoreSuccess(notifConfig.Slack.WebhookURL, dbConfig.Name, backupFilePath)
	}
	return nil
}

func compressBackupFile(inputPath, outputPath string, compressionLevel int) error {
	if err := utils.CompressFileWithLevel(inputPath, outputPath, compressionLevel); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("compressed file was not created: %s", outputPath)
	}

	return nil
}

func uploadToStorage(filePath string, storageConfig config.StorageConfig, dryRun bool) error {
	switch storageConfig.Type {
	case "local":
		zap.L().Sugar().Infof("Backup file saved locally at: %s", filePath)

		if storageConfig.Retain > 0 {
			if err := enforceRetentionLocal(filePath, storageConfig, dryRun); err != nil {
				return err
			}
		}
	case "s3":
		if err := storage.UploadToS3(filePath, storageConfig.Bucket, storageConfig.Region); err != nil {
			return err
		}
		if storageConfig.Retain > 0 {
			if err := storage.EnforceRetentionS3(storageConfig, dryRun); err != nil {
				return err
			}
		}
	case "gcs":
		if err := storage.UploadToGCS(filePath, storageConfig.Bucket, storageConfig.Project); err != nil {
			return err
		}
		if storageConfig.Retain > 0 {
			if err := storage.EnforceRetentionGCS(storageConfig, dryRun); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported storage type: %s", storageConfig.Type)
	}
	return nil
}

func enforceRetentionLocal(newFilePath string, storageConfig config.StorageConfig, dryRun bool) error {
	files, err := os.ReadDir(storageConfig.Path)
	if err != nil {
		return err
	}

	var backupFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			backupPath := filepath.Join(storageConfig.Path, file.Name())
			backupFiles = append(backupFiles, backupPath)
		}
	}

	sort.Slice(backupFiles, func(i, j int) bool {
		t1, _ := time.Parse("20060102_150405", filepath.Base(backupFiles[i]))
		t2, _ := time.Parse("20060102_150405", filepath.Base(backupFiles[j]))
		return t1.After(t2)
	})

	if dryRun {
		zap.L().Sugar().Infof("Simulated deletion of %d old files", len(backupFiles)-storageConfig.Retain)
		return nil
	}

	if len(backupFiles) > storageConfig.Retain {
		for i := len(backupFiles) - 1; i >= storageConfig.Retain; i-- {
			filePath := backupFiles[i]
			if err := os.Remove(filePath); err != nil {
				return err
			}
			zap.L().Sugar().Infof("Deleted old backup: %s", filePath)
		}
	}

	return nil
}
