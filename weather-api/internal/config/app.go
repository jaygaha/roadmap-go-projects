package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Config stuct for app
type Config struct {
	Port          int    `env:"PORT" envDefault:"8800"`
	WeatherAPIKey string `env:"WEATHER_API_KEY,required"`
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
}

// LoadConfig loads config from env
func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	cfg := &Config{}
	if err = env.Parse(cfg); err != nil {
		log.Printf("failed to parse config: %v", err)
		return nil, err
	}

	// Check if required fields are empty
	if cfg.WeatherAPIKey == "" {
		log.Println("WEATHER_API_KEY is required")
		return nil, err
	}

	return cfg, nil
}
