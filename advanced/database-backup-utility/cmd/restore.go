package cmd

import (
	"fmt"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/backup"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var backupFile string
var tablesToRestore string

var restoreCmd = &cobra.Command{
	Use:   "restore <database-name>",
	Short: "Restore a database from backup",
	Long: `Restore a database from a backup file.

You must specify the database name and either provide a backup file path
or use the latest backup from the configured storage.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		databaseName := args[0]

		config, err := config.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}

		dbConfig := findDatabaseConfig(config, databaseName)
		if dbConfig == nil {
			return fmt.Errorf("database not found in config: %s", databaseName)
		}

		var backupFilePath string
		if backupFile != "" {
			backupFilePath = backupFile
		} else {
			return fmt.Errorf("automatic backup discovery not yet implemented. Please specify --file")
		}

		zap.L().Sugar().Infof("Starting restore of %s from %s", databaseName, backupFilePath)

		var tables []string
		if tablesToRestore != "" {
			for _, t := range strings.Split(tablesToRestore, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					tables = append(tables, t)
				}
			}
		}

		if len(tables) > 0 {
			zap.L().Sugar().Infof("Selective restore: tables/collections: %v", tables)
		}

		if err := backup.RestoreDatabase(*dbConfig, backupFilePath, tables, config.Notification); err != nil {
			return fmt.Errorf("restore failed: %v", err)
		}

		zap.L().Sugar().Infof("Restore completed successfully for %s", databaseName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringVarP(&backupFile, "file", "f", "", "Path to backup file to restore from")
	restoreCmd.Flags().StringVarP(&tablesToRestore, "tables", "t", "", "Comma-separated list of tables/collections for selective restore")
}
