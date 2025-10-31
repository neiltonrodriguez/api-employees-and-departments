package logging

import (
	domainLogging "api-employees-and-departments/internal/domain/logging"

	"go.uber.org/zap"
)

// ZapLogger is an adapter that implements domain.Logger using Zap
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger that implements the domain Logger interface
func NewZapLogger(logger *zap.Logger) domainLogging.Logger {
	return &ZapLogger{logger: logger}
}

// Debug logs a debug-level message
func (z *ZapLogger) Debug(msg string, fields ...domainLogging.Field) {
	z.logger.Debug(msg, convertFields(fields)...)
}

// Info logs an info-level message
func (z *ZapLogger) Info(msg string, fields ...domainLogging.Field) {
	z.logger.Info(msg, convertFields(fields)...)
}

// Warn logs a warning-level message
func (z *ZapLogger) Warn(msg string, fields ...domainLogging.Field) {
	z.logger.Warn(msg, convertFields(fields)...)
}

// Error logs an error-level message
func (z *ZapLogger) Error(msg string, fields ...domainLogging.Field) {
	z.logger.Error(msg, convertFields(fields)...)
}

// With creates a child logger with additional fields
func (z *ZapLogger) With(fields ...domainLogging.Field) domainLogging.Logger {
	return &ZapLogger{
		logger: z.logger.With(convertFields(fields)...),
	}
}

// convertFields converts domain logging fields to Zap fields
func convertFields(fields []domainLogging.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = convertField(field)
	}
	return zapFields
}

// convertField converts a single domain field to a Zap field
func convertField(field domainLogging.Field) zap.Field {
	switch v := field.Value.(type) {
	case string:
		return zap.String(field.Key, v)
	case int:
		return zap.Int(field.Key, v)
	case int64:
		return zap.Int64(field.Key, v)
	case bool:
		return zap.Bool(field.Key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.Any(field.Key, v)
	}
}
