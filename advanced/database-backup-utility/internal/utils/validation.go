package utils

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func TestDatabaseConnection(dbConfig config.DatabaseConfig) error {
	switch dbConfig.Type {
	case "mysql":
		return testMySQLConnection(dbConfig)
	case "postgres":
		return testPostgreSQLConnection(dbConfig)
	case "mongodb":
		return testMongoDBConnection(dbConfig)
	case "sqlite":
		return testSQLiteConnection(dbConfig)
	default:
		return fmt.Errorf("unsupported database type for connection testing: %s", dbConfig.Type)
	}
}

func testMongoDBConnection(dbConfig config.DatabaseConfig) error {
	zap.L().Info("Testing MongoDB connectivity (minimal check)", zap.String("database", dbConfig.Name))
	return nil
}

func testSQLiteConnection(dbConfig config.DatabaseConfig) error {
	if _, err := os.Stat(dbConfig.Host); os.IsNotExist(err) {
		return fmt.Errorf("SQLite database file not found at: %s", dbConfig.Host)
	}
	zap.L().Info("SQLite file check successful", zap.String("path", dbConfig.Host))
	return nil
}

func testMySQLConnection(dbConfig config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=5s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %v", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("MySQL connection test failed: %v", err)
	}

	zap.L().Sugar().Infof("MySQL connection test successful for %s", dbConfig.Name)
	return nil
}

func testPostgreSQLConnection(dbConfig config.DatabaseConfig) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=5",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %v", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("PostgreSQL connection test failed: %v", err)
	}

	zap.L().Sugar().Infof("PostgreSQL connection test successful for %s", dbConfig.Name)
	return nil
}

func ValidateDatabaseConfig(dbConfig config.DatabaseConfig) error {
	if dbConfig.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if dbConfig.Type == "" {
		return fmt.Errorf("database type is required")
	}
	if dbConfig.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if dbConfig.User == "" {
		return fmt.Errorf("database user is required")
	}
	if dbConfig.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if dbConfig.Port == 0 {
		return fmt.Errorf("database port is required")
	}

	switch dbConfig.Type {
	case "mysql", "postgres", "mongodb", "sqlite":
	default:
		return fmt.Errorf("unsupported database type: %s", dbConfig.Type)
	}

	if dbConfig.Type == "sqlite" && dbConfig.Host == "" {
		return fmt.Errorf("sqlite requires host to be set to the database file path")
	}

	if dbConfig.Type == "mongodb" && dbConfig.Port == 0 {
		return fmt.Errorf("mongodb requires port to be set")
	}

	return nil
}
