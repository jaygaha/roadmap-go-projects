package scheduler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/backup"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Start launches the backup scheduler and blocks until SIGTERM or SIGINT
func Start(cfg *config.Config) error {
	if !cfg.Schedule.Enabled {
		return fmt.Errorf("scheduling is disabled in config (schedule.enabled: false)")
	}

	if cfg.Schedule.Cron == "" {
		return fmt.Errorf("no cron expression set in config (schedule.cron)")
	}

	loc := time.UTC
	if cfg.Schedule.TimeZone != "" {
		var err error
		loc, err = time.LoadLocation(cfg.Schedule.TimeZone)
		if err != nil {
			return fmt.Errorf("invalid timezone %q: %w", cfg.Schedule.TimeZone, err)
		}
	}

	c := cron.New(cron.WithLocation(loc))

	jobFunc := func() {
		zap.L().Info("Scheduler: starting scheduled backup run",
			zap.String("cron", cfg.Schedule.Cron),
			zap.Time("triggered_at", time.Now()),
		)

		for _, dbCfg := range cfg.Databases {
			zap.L().Info("Scheduler: backing up database", zap.String("database", dbCfg.Name))
			if err := backup.BackupDatabase(dbCfg, cfg.Backup, cfg.Storage, cfg.Notification); err != nil {
				zap.L().Error("Scheduler: backup failed",
					zap.String("database", dbCfg.Name),
					zap.Error(err),
				)
			}
		}

		zap.L().Info("Scheduler: backup run complete")
	}

	entryID, err := c.AddFunc(cfg.Schedule.Cron, jobFunc)
	if err != nil {
		return fmt.Errorf("failed to register cron job (expression: %q): %w", cfg.Schedule.Cron, err)
	}

	c.Start()

	zap.L().Info("Scheduler started",
		zap.String("cron", cfg.Schedule.Cron),
		zap.String("timezone", loc.String()),
		zap.Int("entry_id", int(entryID)),
		zap.Time("next_run", c.Entry(entryID).Next),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	zap.L().Info("Scheduler: shutting down", zap.String("signal", sig.String()))
	ctx := c.Stop()
	<-ctx.Done()
	zap.L().Info("Scheduler: all jobs finished, exiting")
	return nil
}
