package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZerolog_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "info", false)

	logger.Info("This is an info message")

	assert.Contains(t, buf.String(), "This is an info message", "Expected info message to be logged")
}

func TestZerolog_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "debug", false)

	logger.Debug("This is a debug message")

	assert.Contains(t, buf.String(), "This is a debug message", "Expected debug message to be logged")
}

func TestZerolog_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "warn", false)

	logger.Warn("This is a warning message")

	assert.Contains(t, buf.String(), "This is a warning message", "Expected warning message to be logged")
}

func TestZerolog_Info_Parametrised(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "info", false)
	logger.Info("This is an info message with param: %d", 42)
	assert.Contains(t, buf.String(), "This is an info message with param: 42", "Expected info message with parameter to be logged")
}

func TestZerolog_Unknown_UseInfoInstead(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "unknown", false)

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")

	assert.Contains(t, buf.String(), "This is an info message", "Expected info message to be logged")
	assert.NotContains(t, buf.String(), "This is a debug message", "Expected debug message not to be logged")
}

func TestZerolog_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "error", false)

	logger.Error("This is an error message")

	assert.Contains(t, buf.String(), "This is an error message", "Expected error message to be logged")
}

func TestNoColor(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZerolog(&buf, "info", true)

	logger.Info("This is an info message")

	assert.NotContains(t, buf.String(), "\x1b[", "Expected no color codes in the log output")
}
