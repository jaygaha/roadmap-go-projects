package config

import (
	"os"
)

// Config holds all application-wide settings. We load from env vars
// so secrets never live in source code.
type Config struct {
	JWTSecret     string
	StripeKey     string
	DBPath        string
	Port          string
	AdminEmail    string
	AdminPassword string
}

// Load reads environment variables and populates the Config struct.
func Load() *Config {
	return &Config{
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		StripeKey:     getEnv("STRIPE_SECRET_KEY", ""),
		DBPath:        getEnv("DB_PATH", "ecommerce.db"),
		Port:          getEnv("PORT", "8080"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@jaygaha.com.np"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
