package log

import "fmt"

// Mock is a Mock logger that collects log messages for testing purposes.
type Mock struct {
	Messages []string
}

// NewMock creates and returns a new instance of Mock logger.
func NewMock() Logger {
	return &Mock{Messages: []string{}}
}

// Info logs an informational message with optional arguments.
func (m *Mock) Info(msg string, args ...any) {
	m.Messages = append(m.Messages, "mock info: "+fmt.Sprintf(msg, args...))
}

// Debug logs a debug message with optional arguments.
func (m *Mock) Debug(msg string, args ...any) {
	m.Messages = append(m.Messages, fmt.Sprintf("mock debug: %s", fmt.Sprintf(msg, args...)))
}

// Warn logs a warning message with optional arguments.
func (m *Mock) Warn(msg string, args ...any) {
	m.Messages = append(m.Messages, fmt.Sprintf("mock warn: %s", fmt.Sprintf(msg, args...)))
}

// Error logs an error message with optional arguments.
func (m *Mock) Error(msg string, args ...any) {
	m.Messages = append(m.Messages, fmt.Sprintf("mock error: %s", fmt.Sprintf(msg, args...)))
}
