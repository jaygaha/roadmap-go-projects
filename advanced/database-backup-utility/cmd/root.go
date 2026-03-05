package cmd

import (
	"fmt"
	"os"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "database-backup-utility",
	Short: "A CLI utility for backing up and restoring databases",
	Long: `Database Backup Utility is a comprehensive tool for backing up and restoring
various database management systems including MySQL, PostgreSQL, and more.

Features include:
- Full and incremental backups
- Multiple storage backends (local, S3, GCS)
- Compression and encryption
- Automated retention policies
- Detailed logging and notifications`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.InitLogger(); err != nil {
			return fmt.Errorf("error initializing logger: %v", err)
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		_ = zap.L().Sync()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file")
}
