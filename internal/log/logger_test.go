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

func (m *MockLogger) Info(msg string, _ ...any) {
	m.buf.WriteString("INFO: " + msg)
}

func (m *MockLogger) Debug(msg string, _ ...any) {
	m.buf.WriteString("DEBUG: " + msg)
}

func (m *MockLogger) Warn(msg string, _ ...any) {
	m.buf.WriteString("WARN: " + msg)
}

func (m *MockLogger) Error(msg string, _ ...any) {
	m.buf.WriteString("ERROR: " + msg)
}

func TestSetAndGetLogger(t *testing.T) {
	mock := &MockLogger{}
	Set(mock)

	initialized := Default()

	assert.Equal(t, mock, initialized, "Expected retrieved logger to be the mock logger")
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
		Default()
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

func TestNew_WithColorlessTrue_ReturnsDefaultLogger(t *testing.T) {
	p := "colorless-prefix"
	cl := true
	c := Red

	l := New(p, c, cl)

	require.NotNil(t, l, "Logger should not be nil")
	assert.IsType(t, &prefixed{}, l, "Expected a PrefixedLogger")
	pl, ok := l.(*prefixed)
	require.True(t, ok, "Logger should be of type *PrefixedLogger")
	assert.Equal(t, p, pl.prefix, "Prefix should match the input")
	assert.IsType(t, Default(), pl.original, "Inner logger should be ColoredLogger")
}

func TestNew_WithColorlessFalse_ReturnsColoredLogger(t *testing.T) {
	p := "colorfull-prefix"
	cl := false
	c := Green

	l := New(p, c, cl)

	require.NotNil(t, l, "Logger should not be nil")
	assert.IsType(t, &prefixed{}, l, "Expected a PrefixedLogger")
	pl, ok := l.(*prefixed)
	require.True(t, ok, "Logger should be of type *PrefixedLogger")
	assert.Equal(t, p, pl.prefix, "Prefix should match the input")
	assert.IsType(t, &colored{}, pl.original, "Inner logger should be ColoredLogger")
	clogger, ok := pl.original.(*colored)
	require.True(t, ok, "Inner logger should be of type *ColoredLogger")
	assert.Equal(t, c, clogger.color, "Color should match the input")
}
