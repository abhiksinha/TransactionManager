package logger

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Logger wraps zap.Logger and enriches logs with request context when available.
type Logger struct {
	base *zap.Logger
}

// New creates a new logger instance.
func New(isProduction bool) (*Logger, error) {
	var base *zap.Logger
	var err error
	if isProduction {
		base, err = zap.NewProduction()
	} else {
		base, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return &Logger{base: base}, nil
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.base.Sync()
}

// Fatal logs then exits, without request context.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.base.Fatal(msg, fields...)
}

// FromContext returns a logger instance with the request_id from the context baked in.
func FromContext(ctx context.Context, logger *Logger) *zap.Logger {
	reqID := middleware.GetReqID(ctx)
	if reqID != "" {
		return logger.base.With(zap.String("request_id", reqID))
	}
	return logger.base
}

// Debug logs a message with request_id from context if available.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx, l).Debug(msg, fields...)
}

// Info logs a message with request_id from context if available.
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx, l).Info(msg, fields...)
}

// Warn logs a message with request_id from context if available.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx, l).Warn(msg, fields...)
}

// Error logs a message with request_id from context if available.
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx, l).Error(msg, fields...)
}
