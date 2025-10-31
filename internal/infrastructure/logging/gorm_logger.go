package logging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses Zap for structured logging
type GormLogger struct {
	logger                    *zap.Logger
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// NewGormLogger creates a new GORM logger with Zap
func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		logger:                    GetLogger(),
		slowThreshold:             slowThreshold,
		ignoreRecordNotFoundError: true,
	}
}

// LogMode implements gorm logger.Interface
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info implements gorm logger.Interface
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn implements gorm logger.Interface
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error implements gorm logger.Interface
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

// Trace implements gorm logger.Interface for SQL logging
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	// Add request_id from context if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields = append(fields, zap.String("request_id", id))
		}
	}

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		// Log errors (except record not found if ignored)
		fields = append(fields, zap.Error(err))
		l.logger.Error("Database query error", fields...)

	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		// Log slow queries
		fields = append(fields, zap.Duration("threshold", l.slowThreshold))
		l.logger.Warn("Slow database query detected", fields...)

	default:
		// Log normal queries at debug level
		l.logger.Debug("Database query executed", fields...)
	}
}
