package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server  ServerConfig
	Database DatabaseConfig
	JWT     JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret     string
	ExpireTime int // in hours
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("HOST", "localhost"),
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "movie_reservation"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireTime: parseInt(getEnv("JWT_EXPIRE_HOURS", "24")),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseInt converts string to int with default fallback
func parseInt(s string, defaults ...int) int {
	defaultValue := 0
	if len(defaults) > 0 {
		defaultValue = defaults[0]
	}

	if s == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}
