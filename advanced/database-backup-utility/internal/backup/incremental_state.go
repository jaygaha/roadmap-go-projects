package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupState tracks the last backup position for incremental/differential backups
type BackupState struct {
	DBName         string    `json:"db_name"`
	LastFullBackup time.Time `json:"last_full_backup"`
	LastBackupTime time.Time `json:"last_backup_time"`
	LastBackupType string    `json:"last_backup_type"`

	// PostgreSQL-specific
	PGLastLSN string `json:"pg_last_lsn,omitempty"`

	// MySQL-specific
	MySQLBinlogFile string `json:"mysql_binlog_file,omitempty"`
	MySQLBinlogPos  uint32 `json:"mysql_binlog_pos,omitempty"`
}

func stateDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".dbu")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create state directory %s: %w", dir, err)
	}
	return dir, nil
}

func statePath(dbName string) (string, error) {
	dir, err := stateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("%s_state.json", dbName)), nil
}

func LoadState(dbName string) (*BackupState, error) {
	path, err := statePath(dbName)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &BackupState{DBName: dbName}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read state file %s: %w", path, err)
	}

	var state BackupState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file %s: %w", path, err)
	}
	return &state, nil
}

func SaveState(state *BackupState) error {
	path, err := statePath(state.DBName)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write state file %s: %w", path, err)
	}
	return nil
}

func (s *BackupState) IsFirstBackup() bool {
	return s.LastBackupTime.IsZero()
}

func (s *BackupState) NeedsFullBackup() bool {
	return s.LastFullBackup.IsZero()
}
