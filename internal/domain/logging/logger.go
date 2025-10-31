package logging

// Logger is the interface for structured logging in the domain layer.
// This abstraction allows the domain to be independent of specific logging implementations.
type Logger interface {
	// Debug logs a debug-level message with optional fields
	Debug(msg string, fields ...Field)

	// Info logs an info-level message with optional fields
	Info(msg string, fields ...Field)

	// Warn logs a warning-level message with optional fields
	Warn(msg string, fields ...Field)

	// Error logs an error-level message with optional fields
	Error(msg string, fields ...Field)

	// With creates a child logger with additional fields that will be
	// included in all subsequent log calls
	With(fields ...Field) Logger
}

// Field represents a structured logging field with a key-value pair
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new logging field
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// String creates a string field
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Error creates an error field
func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

// Int creates an int field
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a bool field
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Any creates a field with any value type
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
