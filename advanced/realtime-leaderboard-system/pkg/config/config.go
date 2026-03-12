package config

import (
	"os"
	"strconv"
)

var (
	JWTSecret  = getEnv("JWT_SECRET", "default-secret")
	RedisAddr  = getEnv("REDIS_ADDR", "localhost:6379")
	SQLitePath = getEnv("SQLITE_PATH", "./leaderboard.db")
	Port       = getEnv("PORT", "8080")
	RateLimit  = getIntEnv("RATE_LIMIT_PER_MINUTE", 10)
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	strVal := getEnv(key, "")
	if val, err := strconv.Atoi(strVal); err == nil {
		return val
	}
	return fallback
}
