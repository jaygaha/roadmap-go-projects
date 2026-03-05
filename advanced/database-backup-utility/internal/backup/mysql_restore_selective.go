package backup

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
)

// RestoreMySQLSelective restores only the specified tables from a mysqldump SQL file
func RestoreMySQLSelective(dbConfig config.DatabaseConfig, backupFilePath string, tables []string) error {
	tableSet := make(map[string]bool, len(tables))
	for _, t := range tables {
		tableSet[strings.ToLower(t)] = true
	}

	zap.L().Info("Starting MySQL selective restore",
		zap.String("database", dbConfig.Name),
		zap.Strings("tables", tables),
	)

	extracted, err := extractMySQLTableBlocks(backupFilePath, tableSet)
	if err != nil {
		return fmt.Errorf("failed to extract table blocks from dump: %w", err)
	}

	if len(extracted) == 0 {
		return fmt.Errorf("none of the requested tables (%v) were found in the dump file", tables)
	}

	zap.L().Info("Extracted SQL blocks for selective restore",
		zap.String("database", dbConfig.Name),
		zap.Int("blocks", len(extracted)),
	)

	tmpFile, err := os.CreateTemp("", "dbu_selective_*.sql")
	if err != nil {
		return fmt.Errorf("failed to create temp file for selective restore: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	for _, block := range extracted {
		if _, err := fmt.Fprintln(tmpFile, block); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write SQL block to temp file: %w", err)
		}
	}
	tmpFile.Close()

	sqlData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read temp SQL file: %w", err)
	}

	cmd := exec.Command("mysql",
		"-h", dbConfig.Host,
		"-P", fmt.Sprintf("%d", dbConfig.Port),
		"-u", dbConfig.User,
		dbConfig.Name,
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", dbConfig.Password))
	cmd.Stdin = bytes.NewReader(sqlData)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql selective restore failed: %w\nstderr: %s", err, stderr.String())
	}

	zap.L().Info("MySQL selective restore completed",
		zap.String("database", dbConfig.Name),
		zap.Strings("tables", tables),
	)
	return nil
}

func extractMySQLTableBlocks(filePath string, tableSet map[string]bool) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open dump file: %w", err)
	}
	defer f.Close()

	var blocks []string
	var currentBlock strings.Builder
	inBlock := false
	currentTable := ""

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 10*1024*1024)
	scanner.Buffer(buf, 64*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "drop table if exists") {
			tbl := extractTableName(line)
			if tableSet[strings.ToLower(tbl)] {
				blocks = append(blocks, line)
			}
			continue
		}

		if strings.HasPrefix(lower, "create table") {
			tbl := extractTableName(line)
			if tableSet[strings.ToLower(tbl)] {
				inBlock = true
				currentTable = tbl
				currentBlock.Reset()
				currentBlock.WriteString(line + "\n")
			}
			continue
		}

		if strings.HasPrefix(lower, "insert into") {
			tbl := extractTableName(line)
			if tableSet[strings.ToLower(tbl)] {
				blocks = append(blocks, line)
			}
			continue
		}

		if inBlock {
			currentBlock.WriteString(line + "\n")
			trimmed := strings.TrimSpace(line)
			if strings.HasSuffix(trimmed, ";") {
				blocks = append(blocks, fmt.Sprintf("-- Table: %s", currentTable))
				blocks = append(blocks, currentBlock.String())
				inBlock = false
				currentBlock.Reset()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning dump file: %w", err)
	}

	return blocks, nil
}

func extractTableName(line string) string {
	line = strings.ReplaceAll(line, "`", "")
	parts := strings.Fields(line)
	for i, p := range parts {
		pl := strings.ToLower(p)
		if pl == "exists" || pl == "table" || pl == "into" {
			if i+1 < len(parts) {
				return strings.Trim(parts[i+1], ",(;")
			}
		}
	}
	return ""
}
