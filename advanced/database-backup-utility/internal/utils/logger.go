package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func InitLogger() error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("backup_%s.log", time.Now().Format("20060102_150405")))

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zapcore.InfoLevel)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	core := zapcore.NewTee(fileCore, consoleCore)

	l := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	zap.ReplaceGlobals(l)
	globalLogger = l

	l.Info("Logger initialized",
		zap.String("log_file", logFile),
		zap.String("version", "1.0.0"),
	)

	return nil
}

func GetLogger() *zap.Logger {
	return globalLogger
}

func LogOperation(operation string, fields ...zap.Field) func() {
	start := time.Now()

	if globalLogger != nil {
		globalLogger.Info("Operation started",
			append([]zap.Field{zap.String("operation", operation)}, fields...)...)
	}

	return func() {
		duration := time.Since(start)
		if globalLogger != nil {
			globalLogger.Info("Operation completed",
				append([]zap.Field{
					zap.String("operation", operation),
					zap.Duration("duration", duration),
				}, fields...)...)
		}
	}
}
