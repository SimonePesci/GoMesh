package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.Logger
}


func NewLogger(production bool) (*Logger, error) {

	var zapLogger *zap.Logger
	var err error

	if production {
		// Production: Uses JSON encoder and logs from INFO level and above
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger, err = config.Build()

	} else {
		// Development: Human-readable encoder and logs from DEBUG level and above
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = config.Build()

	}

	if err != nil {
		return nil, fmt.Errorf("Failed to create logger: %w", err)
	}

	return &Logger{
		zap: zapLogger,
	}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

// Sync flushes any buffered log entries, put it into easy words: 
// it ensures that all log messages are written to the underlying storage
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// With creates a child logger with additional fields
// Useful for adding context that should be included in all subsequent logs
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		zap: l.zap.With(fields...),
	}
}


// Helper functions for common field types
// These make it easier to add structured data to logs

// String creates a string field
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int creates an integer field
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 creates an int64 field
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Float64 creates a float64 field
func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// Error creates an error field
func Error(err error) zap.Field {
	return zap.Error(err)
}

// Any creates a field from any type (uses reflection, slower)
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
