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

// BackupPostgreSQLIncremental performs an incremental PostgreSQL backup
func BackupPostgreSQLIncremental(dbConfig config.DatabaseConfig, filePath string, state *BackupState) (*BackupState, error) {
	if state.NeedsFullBackup() {
		zap.L().Info("No prior full backup found — falling back to full backup for PostgreSQL incremental",
			zap.String("database", dbConfig.Name),
		)
		if err := BackupPostgreSQL(dbConfig, filePath); err != nil {
			return nil, err
		}
		lsn, err := getPGCurrentLSN(dbConfig)
		if err != nil {
			zap.L().Warn("Could not capture LSN; incremental tracking will be time-based",
				zap.String("database", dbConfig.Name),
				zap.Error(err),
			)
		}
		return &BackupState{
			DBName:         dbConfig.Name,
			LastFullBackup: time.Now(),
			LastBackupTime: time.Now(),
			LastBackupType: "full",
			PGLastLSN:      lsn,
		}, nil
	}

	zap.L().Info("Starting PostgreSQL incremental backup",
		zap.String("database", dbConfig.Name),
		zap.Time("since", state.LastBackupTime),
		zap.String("last_lsn", state.PGLastLSN),
	)

	modifiedTables, err := getPGTablesModifiedSince(dbConfig, state.LastBackupTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query modified tables: %w", err)
	}

	if len(modifiedTables) == 0 {
		zap.L().Info("No tables modified since last backup — creating empty incremental backup",
			zap.String("database", dbConfig.Name),
		)
		if err := os.WriteFile(filePath, []byte("-- No changes since last backup\n"), 0644); err != nil {
			return nil, fmt.Errorf("failed to write empty incremental file: %w", err)
		}
	} else {
		args := buildPGIncrementalArgs(dbConfig, filePath, modifiedTables)
		cmd := exec.Command("pg_dump", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("pg_dump incremental failed: %w\nstderr: %s", err, stderr.String())
		}

		zap.L().Info("PostgreSQL incremental backup completed",
			zap.String("database", dbConfig.Name),
			zap.String("backup_file", filePath),
			zap.Int("tables_backed_up", len(modifiedTables)),
			zap.Strings("tables", modifiedTables),
		)
	}

	lsn, _ := getPGCurrentLSN(dbConfig)

	return &BackupState{
		DBName:         dbConfig.Name,
		LastFullBackup: state.LastFullBackup,
		LastBackupTime: time.Now(),
		LastBackupType: "incremental",
		PGLastLSN:      lsn,
	}, nil
}

// BackupPostgreSQLDifferential performs a differential PostgreSQL backup
func BackupPostgreSQLDifferential(dbConfig config.DatabaseConfig, filePath string, state *BackupState) (*BackupState, error) {
	if state.NeedsFullBackup() {
		zap.L().Info("No prior full backup — falling back to full for PostgreSQL differential",
			zap.String("database", dbConfig.Name),
		)
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

	zap.L().Info("Starting PostgreSQL differential backup (since last full)",
		zap.String("database", dbConfig.Name),
		zap.Time("since_full", state.LastFullBackup),
	)

	modifiedTables, err := getPGTablesModifiedSince(dbConfig, state.LastFullBackup)
	if err != nil {
		return nil, fmt.Errorf("failed to query modified tables for differential: %w", err)
	}

	if len(modifiedTables) == 0 {
		zap.L().Info("No tables modified since last full backup — creating empty differential backup",
			zap.String("database", dbConfig.Name),
		)
		if err := os.WriteFile(filePath, []byte("-- No changes since last full backup\n"), 0644); err != nil {
			return nil, fmt.Errorf("failed to write empty differential file: %w", err)
		}
	} else {
		args := buildPGIncrementalArgs(dbConfig, filePath, modifiedTables)
		cmd := exec.Command("pg_dump", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("pg_dump differential failed: %w\nstderr: %s", err, stderr.String())
		}

		zap.L().Info("PostgreSQL differential backup completed",
			zap.String("database", dbConfig.Name),
			zap.String("backup_file", filePath),
			zap.Int("tables_backed_up", len(modifiedTables)),
		)
	}

	lsn, _ := getPGCurrentLSN(dbConfig)

	return &BackupState{
		DBName:         dbConfig.Name,
		LastFullBackup: state.LastFullBackup,
		LastBackupTime: time.Now(),
		LastBackupType: "differential",
		PGLastLSN:      lsn,
	}, nil
}

func getPGTablesModifiedSince(dbConfig config.DatabaseConfig, since time.Time) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT schemaname || '.' || relname
		FROM pg_stat_user_tables
		WHERE (last_autovacuum > '%s'::timestamptz
		   OR last_vacuum    > '%s'::timestamptz
		   OR last_autoanalyze > '%s'::timestamptz
		   OR last_analyze   > '%s'::timestamptz
		   OR n_mod_since_analyze > 0)
		ORDER BY relname;
	`, since.UTC().Format(time.RFC3339),
		since.UTC().Format(time.RFC3339),
		since.UTC().Format(time.RFC3339),
		since.UTC().Format(time.RFC3339),
	)

	args := []string{
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"-w",
		"-t",
		"-A",
		"-c", query,
	}

	cmd := exec.Command("psql", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("psql query failed: %w — %s", err, stderr.String())
	}

	var tables []string
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		t := strings.TrimSpace(line)
		if t != "" {
			tables = append(tables, t)
		}
	}
	return tables, nil
}

func buildPGIncrementalArgs(dbConfig config.DatabaseConfig, filePath string, tables []string) []string {
	args := []string{
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"-f", filePath,
		"-w",
		"--section=pre-data",
		"--section=data",
	}
	for _, t := range tables {
		args = append(args, fmt.Sprintf("--table=%s", t))
	}
	return args
}

func getPGCurrentLSN(dbConfig config.DatabaseConfig) (string, error) {
	args := []string{
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"-w",
		"-t",
		"-A",
		"-c", "SELECT pg_current_wal_lsn();",
	}

	cmd := exec.Command("psql", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not query LSN: %w — %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
