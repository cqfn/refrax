// Package log provides utilities for structured logging in applications.
package log

import "os"

// Logger is an interface that defines methods for structured logging.
type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

var single = NewZerolog(os.Stdout, "info")

// Info logs an informational message using the default logger.
func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

// Debug logs a debug message using the default logger.
func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

// Error logs an error message using the default logger.
func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

// Set sets the provided logger as the default logger.
func Set(logger Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}
	single = logger
}

// Default returns the currently set default logger.
func Default() Logger {
	if single == nil {
		panic("logger is not set")
	}
	return single
}
