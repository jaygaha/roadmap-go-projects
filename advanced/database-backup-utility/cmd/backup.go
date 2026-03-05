package cmd

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/backup"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var backupType string

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Perform database backups",
	Long: `Backup one or more databases according to the configuration.

This command will backup all databases configured in the config file,
or specific databases if names are provided as arguments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig(cfgFile)
		if err != nil {
			return utils.HandleError(err, "Failed to load configuration")
		}

		if backupType != "" {
			config.Backup.Type = backupType
		}

		defer utils.LogOperation("backup_process",
			zap.String("config_file", cfgFile),
			zap.Int("database_count", len(config.Databases)),
		)()

		if len(args) > 0 {
			for _, dbName := range args {
				if dbConfig := findDatabaseConfig(config, dbName); dbConfig != nil {
					if err := backup.BackupDatabase(*dbConfig, config.Backup, config.Storage, config.Notification); err != nil {
						return fmt.Errorf("backup failed for %s: %v", dbName, err)
					}
				} else {
					return fmt.Errorf("database not found in config: %s", dbName)
				}
			}
		} else {
			for _, dbConfig := range config.Databases {
				if err := backup.BackupDatabase(dbConfig, config.Backup, config.Storage, config.Notification); err != nil {
					return fmt.Errorf("backup failed for %s: %v", dbConfig.Name, err)
				}
			}
		}

		zap.L().Info("Backup process completed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&backupType, "type", "t", "", "Backup type: full, incremental, or differential (overrides config)")
}

func findDatabaseConfig(cfg *config.Config, name string) *config.DatabaseConfig {
	for _, db := range cfg.Databases {
		if db.Name == name {
			return &db
		}
	}
	return nil
}
