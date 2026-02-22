package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	MongoURI            string
	AWSEndpoint         string
	AWSAccessKey        string
	AWSSecretKey        string
	AWSRegion           string
	S3Bucket            string
	Port                string
	JWTSecret           string
	JWTIssuer           string
	JWTExpirationMinute int
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig() *Config {
	godotenv.Load()

	jwtExpirationMinute := getEnvInt("JWT_EXPIRATION_MINUTE", 60)

	return &Config{
		MongoURI:            getEnv("MONGODB_URI", "mongodb://root:root%40123@mongodb:27017/image_processing?authSource=admin"),
		AWSEndpoint:         getEnv("S3_ENDPOINT_URL", "http://localhost:4566"),
		AWSAccessKey:        getEnv("AWS_ACCESS_KEY_ID", "localstack"),
		AWSSecretKey:        getEnv("AWS_SECRET_ACCESS_KEY", "localstack"),
		AWSRegion:           getEnv("AWS_REGION", "us-east-1"),
		S3Bucket:            getEnv("AWS_BUCKET", "img-bucket"),
		Port:                getEnv("APP_PORT", "8800"),
		JWTSecret:           getEnv("JWT_SECRET", "supersecretjwtkey"),
		JWTIssuer:           getEnv("JWT_ISSUER", "image-processing-service"),
		JWTExpirationMinute: jwtExpirationMinute,
	}
}

// getEnv returns the environment variable value or the default value if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		n, err := strconv.Atoi(value)
		if err == nil {
			return n
		}
	}

	return defaultValue
}
