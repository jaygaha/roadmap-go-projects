package cmd

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test database connections",
	Long: `Test database connections to ensure configurations are correct.

This command will test the connection to all configured databases
or specific databases if names are provided as arguments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}

		zap.L().Info("Starting database connection tests...")

		if len(args) > 0 {
			for _, dbName := range args {
				if dbConfig := findDatabaseConfig(config, dbName); dbConfig != nil {
					if err := testSingleDatabase(*dbConfig); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("database not found in config: %s", dbName)
				}
			}
		} else {
			for _, dbConfig := range config.Databases {
				if err := testSingleDatabase(dbConfig); err != nil {
					return err
				}
			}
		}

		zap.L().Info("All database connection tests completed successfully!")
		return nil
	},
}

func testSingleDatabase(dbConfig config.DatabaseConfig) error {
	if err := utils.ValidateDatabaseConfig(dbConfig); err != nil {
		return fmt.Errorf("invalid configuration for %s: %v", dbConfig.Name, err)
	}

	zap.L().Sugar().Infof("Testing connection to %s (%s)...", dbConfig.Name, dbConfig.Type)

	if err := utils.TestDatabaseConnection(dbConfig); err != nil {
		return fmt.Errorf("connection test failed for %s: %v", dbConfig.Name, err)
	}

	zap.L().Sugar().Infof("✓ Connection to %s successful", dbConfig.Name)
	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
