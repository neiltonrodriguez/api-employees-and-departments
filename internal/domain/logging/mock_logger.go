package logging

// MockLogger is a test implementation of Logger interface for unit testing
// Usage example in tests:
//
//	mock := logging.NewMockLogger()
//	service := employee.NewService(repo, mock)
//	// ... test service methods
//	// ... assert mock.Logs contains expected entries
type MockLogger struct {
	Logs []LogEntry
}

// LogEntry represents a single log entry for testing
type LogEntry struct {
	Level   string
	Message string
	Fields  []Field
}

// NewMockLogger creates a new mock logger for testing
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Logs: make([]LogEntry, 0),
	}
}

func (m *MockLogger) Debug(msg string, fields ...Field) {
	m.Logs = append(m.Logs, LogEntry{Level: "DEBUG", Message: msg, Fields: fields})
}

func (m *MockLogger) Info(msg string, fields ...Field) {
	m.Logs = append(m.Logs, LogEntry{Level: "INFO", Message: msg, Fields: fields})
}

func (m *MockLogger) Warn(msg string, fields ...Field) {
	m.Logs = append(m.Logs, LogEntry{Level: "WARN", Message: msg, Fields: fields})
}

func (m *MockLogger) Error(msg string, fields ...Field) {
	m.Logs = append(m.Logs, LogEntry{Level: "ERROR", Message: msg, Fields: fields})
}

func (m *MockLogger) With(fields ...Field) Logger {
	// For simplicity, return the same mock logger
	// In a more sophisticated implementation, could create a child with prefixed fields
	return m
}

// Reset clears all logged entries (useful between test cases)
func (m *MockLogger) Reset() {
	m.Logs = make([]LogEntry, 0)
}

// GetLastLog returns the most recent log entry
func (m *MockLogger) GetLastLog() *LogEntry {
	if len(m.Logs) == 0 {
		return nil
	}
	return &m.Logs[len(m.Logs)-1]
}

// CountByLevel counts log entries by level
func (m *MockLogger) CountByLevel(level string) int {
	count := 0
	for _, entry := range m.Logs {
		if entry.Level == level {
			count++
		}
	}
	return count
}

// HasMessage checks if any log contains the given message
func (m *MockLogger) HasMessage(message string) bool {
	for _, entry := range m.Logs {
		if entry.Message == message {
			return true
		}
	}
	return false
}
