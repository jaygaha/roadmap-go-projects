package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap logger
func New() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()
	return logger
}

// Helper functions for consistent error logging
// Error logs an error with the given message and fields
func Error(err error) zap.Field {
	return zap.Error(err)
}

// String logs a string with the given key and value
func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

// Int logs an integer with the given key and value
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Any logs any value with the given key and value
func Any(key string, value any) zap.Field {
	return zap.Any(key, value)
}
