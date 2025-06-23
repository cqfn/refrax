package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock_Info(t *testing.T) {
	mockLogger := NewMock()

	mockLogger.Info("This is an info message")

	assert.Len(t, mockLogger.Messages, 1, "Expected one message to be logged")
	assert.Contains(t, mockLogger.Messages[0], "mock info: This is an info message", "Expected info message to be logged")
}

func TestMock_Debug(t *testing.T) {
	mockLogger := NewMock()

	mockLogger.Debug("This is a debug message")

	assert.Len(t, mockLogger.Messages, 1, "Expected one message to be logged")
	assert.Contains(t, mockLogger.Messages[0], "mock debug: This is a debug message", "Expected debug message to be logged")
}

func TestMock_Warn(t *testing.T) {
	mockLogger := NewMock()

	mockLogger.Warn("This is a warning message")

	assert.Len(t, mockLogger.Messages, 1, "Expected one message to be logged")
	assert.Contains(t, mockLogger.Messages[0], "mock warn: This is a warning message", "Expected warning message to be logged")
}
