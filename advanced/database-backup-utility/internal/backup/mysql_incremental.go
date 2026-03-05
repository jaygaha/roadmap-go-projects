package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// BackupMySQLIncremental performs an incremental MySQL backup using mysqlbinlog
func BackupMySQLIncremental(dbConfig config.DatabaseConfig, filePath string, state *BackupState) (*BackupState, error) {
	if state.NeedsFullBackup() {
		zap.L().Info("No prior full backup found — falling back to full backup for MySQL incremental",
			zap.String("database", dbConfig.Name),
		)
		if err := BackupMySQL(dbConfig, filePath); err != nil {
			return nil, err
		}
		newState := &BackupState{
			DBName:         dbConfig.Name,
			LastFullBackup: time.Now(),
			LastBackupTime: time.Now(),
			LastBackupType: "full",
		}
		file, pos, err := getMySQLBinlogPosition(dbConfig)
		if err != nil {
			zap.L().Warn("Could not capture binlog position; incremental backups may not work",
				zap.String("database", dbConfig.Name),
				zap.Error(err),
			)
		} else {
			newState.MySQLBinlogFile = file
			newState.MySQLBinlogPos = pos
		}
		return newState, nil
	}

	if state.MySQLBinlogFile == "" {
		return nil, fmt.Errorf("no binlog position recorded for %s; run a full backup first", dbConfig.Name)
	}

	zap.L().Info("Starting MySQL incremental backup",
		zap.String("database", dbConfig.Name),
		zap.String("binlog_file", state.MySQLBinlogFile),
		zap.Uint32("binlog_pos", state.MySQLBinlogPos),
	)

	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create incremental backup file: %w", err)
	}
	defer outFile.Close()

	args := []string{
		"--read-from-remote-server",
		fmt.Sprintf("--host=%s", dbConfig.Host),
		fmt.Sprintf("--port=%d", dbConfig.Port),
		fmt.Sprintf("--user=%s", dbConfig.User),
		fmt.Sprintf("--start-position=%d", state.MySQLBinlogPos),
		"--to-last-log",
		"--result-file=/dev/stdout",
		state.MySQLBinlogFile,
	}

	cmd := exec.Command("mysqlbinlog", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", dbConfig.Password))
	cmd.Stdout = outFile

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mysqlbinlog incremental backup failed: %w\nstderr: %s", err, stderr.String())
	}

	newFile, newPos, err := getMySQLBinlogPosition(dbConfig)
	if err != nil {
		zap.L().Warn("Could not update binlog position after incremental backup",
			zap.String("database", dbConfig.Name),
			zap.Error(err),
		)
		newFile = state.MySQLBinlogFile
		newPos = state.MySQLBinlogPos
	}

	newState := &BackupState{
		DBName:          dbConfig.Name,
		LastFullBackup:  state.LastFullBackup,
		LastBackupTime:  time.Now(),
		LastBackupType:  "incremental",
		MySQLBinlogFile: newFile,
		MySQLBinlogPos:  newPos,
	}

	zap.L().Info("MySQL incremental backup completed",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", filePath),
		zap.String("new_binlog_file", newFile),
		zap.Uint32("new_binlog_pos", newPos),
	)

	return newState, nil
}

// BackupMySQLDifferential performs a differential MySQL backup
func BackupMySQLDifferential(dbConfig config.DatabaseConfig, filePath string, state *BackupState) (*BackupState, error) {
	if state.NeedsFullBackup() {
		zap.L().Info("No prior full backup found — falling back to full backup for MySQL differential",
			zap.String("database", dbConfig.Name),
		)
		if err := BackupMySQL(dbConfig, filePath); err != nil {
			return nil, err
		}
		newState := &BackupState{
			DBName:         dbConfig.Name,
			LastFullBackup: time.Now(),
			LastBackupTime: time.Now(),
			LastBackupType: "full",
		}
		file, pos, err := getMySQLBinlogPosition(dbConfig)
		if err == nil {
			newState.MySQLBinlogFile = file
			newState.MySQLBinlogPos = pos
		}
		return newState, nil
	}

	zap.L().Info("Starting MySQL differential backup (since last full)",
		zap.String("database", dbConfig.Name),
		zap.Time("since_full", state.LastFullBackup),
	)

	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create differential backup file: %w", err)
	}
	defer outFile.Close()

	args := []string{
		"--read-from-remote-server",
		fmt.Sprintf("--host=%s", dbConfig.Host),
		fmt.Sprintf("--port=%d", dbConfig.Port),
		fmt.Sprintf("--user=%s", dbConfig.User),
		fmt.Sprintf("--start-position=%d", state.MySQLBinlogPos),
		"--to-last-log",
		"--result-file=/dev/stdout",
		state.MySQLBinlogFile,
	}

	cmd := exec.Command("mysqlbinlog", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", dbConfig.Password))
	cmd.Stdout = outFile

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mysqlbinlog differential backup failed: %w\nstderr: %s", err, stderr.String())
	}

	newState := &BackupState{
		DBName:          dbConfig.Name,
		LastFullBackup:  state.LastFullBackup,
		LastBackupTime:  time.Now(),
		LastBackupType:  "differential",
		MySQLBinlogFile: state.MySQLBinlogFile,
		MySQLBinlogPos:  state.MySQLBinlogPos,
	}

	zap.L().Info("MySQL differential backup completed",
		zap.String("database", dbConfig.Name),
		zap.String("backup_file", filePath),
	)

	return newState, nil
}

func getMySQLBinlogPosition(dbConfig config.DatabaseConfig) (string, uint32, error) {
	args := []string{
		"-h", dbConfig.Host,
		"-P", fmt.Sprintf("%d", dbConfig.Port),
		"-u", dbConfig.User,
		"--execute=SHOW MASTER STATUS",
		"--silent",
		"--skip-column-names",
	}

	cmd := exec.Command("mysql", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", dbConfig.Password))

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", 0, fmt.Errorf("SHOW MASTER STATUS failed: %w — %s", err, stderr.String())
	}

	line := strings.TrimSpace(out.String())
	if line == "" {
		return "", 0, fmt.Errorf("binary logging may not be enabled on the MySQL server")
	}

	parts := strings.Split(line, "\t")
	if len(parts) < 2 {
		return "", 0, fmt.Errorf("unexpected SHOW MASTER STATUS output: %q", line)
	}

	var pos uint32
	if _, err := fmt.Sscanf(parts[1], "%d", &pos); err != nil {
		return "", 0, fmt.Errorf("could not parse binlog position %q: %w", parts[1], err)
	}

	return parts[0], pos, nil
}
