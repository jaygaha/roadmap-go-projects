package repository

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/migrations"
)

// RunMigrations applies all pending database migrations in order.
// It tracks applied migrations in a schema_migrations table so each
// migration runs exactly once.
func (db *DB) RunMigrations(ctx context.Context) error {
	// Ensure the tracking table exists
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	// Collect and sort migration files by name (lexicographic = numeric order)
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version := strings.TrimSuffix(entry.Name(), ".sql")

		// Skip already-applied migrations
		var applied bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`,
			version,
		).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if applied {
			continue
		}

		content, err := migrations.FS.ReadFile(entry.Name())
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", version, err)
		}

		upSQL := extractUpSQL(string(content))

		// Apply in a transaction so partial failures roll back cleanly
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin transaction for migration %s: %w", version, err)
		}

		if strings.TrimSpace(upSQL) != "" {
			if _, err := tx.ExecContext(ctx, upSQL); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("apply migration %s: %w", version, err)
			}
		}

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", version, err)
		}

		log.Printf("[migration] applied: %s", version)
	}

	return nil
}

// extractUpSQL returns the SQL between the "-- UP" and "-- DOWN" markers.
// If no markers are present the entire file content is returned.
func extractUpSQL(content string) string {
	const (
		upMarker   = "-- UP"
		downMarker = "-- DOWN"
	)

	upIdx := strings.Index(content, upMarker)
	if upIdx == -1 {
		return strings.TrimSpace(content)
	}

	upStart := upIdx + len(upMarker)
	rest := content[upStart:]

	downIdx := strings.Index(rest, downMarker)
	if downIdx == -1 {
		return strings.TrimSpace(rest)
	}

	return strings.TrimSpace(rest[:downIdx])
}
