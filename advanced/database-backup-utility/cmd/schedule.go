package cmd

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/scheduler"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Run backups on an automated schedule",
	Long: `Start a long-running scheduler that automatically backs up configured databases
on a cron schedule. Reads the cron expression and timezone from config.yaml
(schedule.cron and schedule.timezone), or override them with flags.

Examples:
  # Use schedule settings from config.yaml
  dbu schedule --config config.yaml

  # Override cron expression at runtime
  dbu schedule --cron "0 2 * * *"

  # Override and set timezone
  dbu schedule --cron "0 2 * * *" --timezone "Asia/Tokyo"

The scheduler blocks until it receives SIGINT or SIGTERM, then waits for any
in-progress backup to finish before exiting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %v", err)
		}

		if cronExpr, _ := cmd.Flags().GetString("cron"); cronExpr != "" {
			cfg.Schedule.Cron = cronExpr
			cfg.Schedule.Enabled = true
		}
		if tz, _ := cmd.Flags().GetString("timezone"); tz != "" {
			cfg.Schedule.TimeZone = tz
		}

		if !cfg.Schedule.Enabled {
			cfg.Schedule.Enabled = true
		}

		zap.L().Info("Starting backup scheduler",
			zap.String("cron", cfg.Schedule.Cron),
			zap.String("timezone", cfg.Schedule.TimeZone),
			zap.Int("databases", len(cfg.Databases)),
		)

		return scheduler.Start(cfg)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().String("cron", "", `Cron expression override, e.g. "0 2 * * *" (daily at 02:00)`)
	scheduleCmd.Flags().String("timezone", "", `Timezone override, e.g. "Asia/Tokyo" (default UTC)`)
}
