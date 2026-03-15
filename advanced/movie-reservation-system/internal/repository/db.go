package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DB wraps the database connection with context support
type DB struct {
	*sql.DB
	config *DBConfig
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// New creates a new database connection
func New(config *DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db, config: config}, nil
}

// ExecContext executes a query with context
func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.DB.ExecContext(ctx, query, args...)
}

// QueryContext executes a query with context and returns rows
func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query with context and returns a single row
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRowContext(ctx, query, args...)
}

// BeginTx starts a transaction with context
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.DB.BeginTx(ctx, opts)
}

// Config returns the database configuration
func (d *DB) Config() *DBConfig {
	return d.config
}
