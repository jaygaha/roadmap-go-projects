package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// RestorePostgreSQLSelective restores only the specified tables from a PostgreSQL backup file
func RestorePostgreSQLSelective(dbConfig config.DatabaseConfig, backupFilePath string, tables []string) error {
	zap.L().Info("Starting PostgreSQL selective restore",
		zap.String("database", dbConfig.Name),
		zap.Strings("tables", tables),
	)

	if strings.HasSuffix(backupFilePath, ".pgdump") {
		return restorePostgreSQLSelectiveCustomFmt(dbConfig, backupFilePath, tables)
	}
	return restorePostgreSQLSelectivePlainSQL(dbConfig, backupFilePath, tables)
}

func restorePostgreSQLSelectiveCustomFmt(dbConfig config.DatabaseConfig, backupFilePath string, tables []string) error {
	args := []string{
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"--no-owner",
		"--no-privileges",
		"--single-transaction",
	}

	for _, t := range tables {
		parts := strings.SplitN(t, ".", 2)
		tableName := parts[len(parts)-1]
		args = append(args, fmt.Sprintf("--table=%s", tableName))
	}

	args = append(args, backupFilePath)

	cmd := exec.Command("pg_restore", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_restore selective failed: %w\nstderr: %s", err, stderr.String())
	}

	zap.L().Info("PostgreSQL selective restore (custom format) completed",
		zap.String("database", dbConfig.Name),
		zap.Strings("tables", tables),
	)
	return nil
}

func restorePostgreSQLSelectivePlainSQL(dbConfig config.DatabaseConfig, backupFilePath string, tables []string) error {
	tableSet := make(map[string]bool, len(tables))
	for _, t := range tables {
		parts := strings.SplitN(t, ".", 2)
		tableSet[strings.ToLower(parts[len(parts)-1])] = true
	}

	sqlData, err := os.ReadFile(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL dump file: %w", err)
	}

	filtered := filterPostgreSQLStatements(string(sqlData), tableSet)
	if strings.TrimSpace(filtered) == "" {
		return fmt.Errorf("no SQL statements found for requested tables %v in dump file", tables)
	}

	cmd := exec.Command("psql",
		"-h", dbConfig.Host,
		"-p", fmt.Sprintf("%d", dbConfig.Port),
		"-U", dbConfig.User,
		"-d", dbConfig.Name,
		"-w",
		"--single-transaction",
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))
	cmd.Stdin = strings.NewReader(filtered)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql selective restore failed: %w\nstderr: %s", err, stderr.String())
	}

	zap.L().Info("PostgreSQL selective restore (plain SQL) completed",
		zap.String("database", dbConfig.Name),
		zap.Strings("tables", tables),
	)
	return nil
}

func filterPostgreSQLStatements(sql string, tableSet map[string]bool) string {
	var result strings.Builder
	lines := strings.Split(sql, "\n")

	inBlock := false
	blockDepth := 0

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))

		if strings.HasPrefix(lower, "create table") {
			tbl := extractPGTableName(line)
			if tableSet[strings.ToLower(tbl)] {
				inBlock = true
				blockDepth = 0
				result.WriteString(line + "\n")
				continue
			}
		}

		if strings.HasPrefix(lower, "copy ") {
			tbl := extractPGTableName(line)
			if tableSet[strings.ToLower(tbl)] {
				inBlock = true
				blockDepth = 0
				result.WriteString(line + "\n")
				continue
			} else {
				inBlock = false
			}
		}

		if inBlock {
			result.WriteString(line + "\n")

			blockDepth += strings.Count(line, "(") - strings.Count(line, ")")

			trimmed := strings.TrimSpace(line)
			if blockDepth <= 0 && (strings.HasSuffix(trimmed, ");") || strings.HasSuffix(trimmed, ");")) {
				inBlock = false
			}
			if trimmed == "\\." {
				inBlock = false
			}
		}
	}

	return result.String()
}

func extractPGTableName(line string) string {
	line = strings.ReplaceAll(line, "\"", "")
	parts := strings.Fields(line)
	for i, p := range parts {
		pl := strings.ToLower(p)
		if pl == "table" || pl == "copy" {
			if i+1 < len(parts) {
				name := strings.Trim(parts[i+1], "(,;")
				segments := strings.SplitN(name, ".", 2)
				return segments[len(segments)-1]
			}
		}
	}
	return ""
}
