package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	AppPort   string
	DBDriver  string
	DBName    string
	JWTSecret string
	JWTEXP    int
}

// LoadConfig loads configuration from .env and environment variables
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil
	}

	jwtExpInt, err := strconv.Atoi(getEnvVariable("JWT_EXP", "1"))
	if err != nil {
		log.Fatal("Error converting JWT_EXP to integer")
		return nil
	}

	// Create a new Config instance
	config := &Config{
		AppPort:   getEnvVariable("APP_PORT", "8800"),
		DBDriver:  getEnvVariable("DB_DRIVER", "sqlite3"),
		DBName:    getEnvVariable("DB_NAME", "expense_tracker.db"),
		JWTSecret: getEnvVariable("JWT_SECRET", ""),
		JWTEXP:    jwtExpInt,
	}

	return config
}

// getEnvVariable gets the value of an environment variable
func getEnvVariable(key, defaultValue string) string {
	// Check if the environment variable exists
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
