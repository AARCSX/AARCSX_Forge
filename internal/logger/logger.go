package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger to provide a consistent interface.
type Logger struct {
	*zap.Logger
}

// New creates a new logger based on the provided level and service name.
// level: debug, info, warn, error, dpanic, panic, fatal.
// serviceName: name of the service for logging context.
func New(level string, serviceName string) (*Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "dpanic":
		zapLevel = zapcore.DPanicLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.InitialFields = map[string]interface{}{
		"service": serviceName,
	}
	// Disable sampler to keep it simple for now.
	config.Sampling = nil

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger}, nil
}

// Sugar returns a SugaredLogger for less verbose logging.
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.Logger.Sugar()
}

// With creates a child logger with additional fields.
// This is useful for adding context-specific fields like request_id, trace_id, etc.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

// SugarWith returns a SugaredLogger with additional fields.
func (l *Logger) SugarWith(fields ...zap.Field) *zap.SugaredLogger {
	return l.Logger.With(fields...).Sugar()
}

// Info logs an informational message with context and fields.
func (l *Logger) Info(ctx context.Context, msg string, fields map[string]any) {
	if l == nil { return }
	zapFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}
	l.Logger.Sugar().Infow(msg, zapFields...)
}

// Error logs an error message with context, associated error, and fields.
func (l *Logger) Error(ctx context.Context, msg string, err error, fields map[string]any) {
	if l == nil { return }
	zapFields := make([]interface{}, 0, len(fields)*2+2)
	zapFields = append(zapFields, "error", err)
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}
	l.Logger.Sugar().Errorw(msg, zapFields...)
}