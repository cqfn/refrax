package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	buf bytes.Buffer
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.buf.WriteString("INFO: " + msg)
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.buf.WriteString("DEBUG: " + msg)
}

func (m *MockLogger) Warn(msg string, args ...any) {
	m.buf.WriteString("WARN: " + msg)
}

func (m *MockLogger) Error(msg string, args ...any) {
	m.buf.WriteString("ERROR: " + msg)
}

func TestSetAndGetLogger(t *testing.T) {
	mock := &MockLogger{}
	Set(mock)

	initialised := Get()

	assert.Equal(t, mock, initialised, "Expected retrieved logger to be the mock logger")
}

func TestSetLoggerNilPanics(t *testing.T) {
	assert.Panics(t, func() {
		Set(nil)
	}, "Expected panic when setting logger to nil")
}

func TestGetLoggerNotSetPanics(t *testing.T) {
	original := single
	defer func() {
		single = original
	}()
	single = nil
	require.Panics(t, func() {
		Get()
	}, "Expected panic when getting logger that is not set")
}

func TestLogger_Info(t *testing.T) {
	mock := &MockLogger{}
	Set(mock)

	Info("This is an info message")

	assert.Contains(t, mock.buf.String(), "INFO: This is an info message", "Expected info message to be logged")
}

func TestLogger_Debug(t *testing.T) {
	mock := &MockLogger{}
	Set(mock)

	Debug("This is a debug message")

	assert.Contains(t, mock.buf.String(), "DEBUG: This is a debug message", "Expected debug message to be logged")
}

func TestLogger_Warn(t *testing.T) {
	mock := &MockLogger{}
	Set(mock)

	Warn("This is a warning message")

	assert.Contains(t, mock.buf.String(), "WARN: This is a warning message", "Expected warning message to be logged")
}

func TestLogger_Error(t *testing.T) {
	mock := &MockLogger{}
	Set(mock)

	Error("This is an error message")

	assert.Contains(t, mock.buf.String(), "ERROR: This is an error message", "Expected error message to be logged")
}
