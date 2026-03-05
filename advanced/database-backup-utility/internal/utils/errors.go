package utils

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// HandleError logs and returns an error with context
func HandleError(err error, message string, fields ...zap.Field) error {
	if err == nil {
		return nil
	}

	LogError(message, err, fields...)
	return fmt.Errorf("%s: %w", message, err)
}

// LogError logs an error with optional context fields
func LogError(message string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err)}, fields...)
	zap.L().Error(message, allFields...)
}

// RetryOperation retries an operation with exponential backoff
func RetryOperation(operation func() error, maxAttempts int, operationName string) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err
		zap.L().Error("Operation failed",
			zap.String("operation", operationName),
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxAttempts),
			zap.Error(err),
		)

		if attempt < maxAttempts {
			sleepDuration := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(sleepDuration)
		}
	}

	return fmt.Errorf("%s failed after %d attempts: %w", operationName, maxAttempts, lastErr)
}
