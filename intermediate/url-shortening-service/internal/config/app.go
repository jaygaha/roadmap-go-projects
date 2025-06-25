package config

import (
	"os"
)

// Config represents application configuration
type config struct {
	MongoDBURI   string
	DatabaseName string
	BaseURL      string
}

// LoadConfig bootstraps the application configuration
func LoadConfig() *config {
	return &config{
		MongoDBURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: "urldb",
		BaseURL:      "http://localhost:8800",
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
