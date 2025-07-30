package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLogger struct {
	infoMessages  []string
	debugMessages []string
	warnMessages  []string
	errorMessages []string
}

func (m *mockLogger) Info(msg string, _ ...any) {
	m.infoMessages = append(m.infoMessages, msg)
}

func (m *mockLogger) Debug(msg string, _ ...any) {
	m.debugMessages = append(m.debugMessages, msg)
}

func (m *mockLogger) Warn(msg string, _ ...any) {
	m.warnMessages = append(m.warnMessages, msg)
}

func (m *mockLogger) Error(msg string, _ ...any) {
	m.errorMessages = append(m.errorMessages, msg)
}

func TestNewColored_Info(t *testing.T) {
	m := &mockLogger{}
	c := NewColored(m, Red)

	require.NotNil(t, c)
	c.Info("test message")

	require.Len(t, m.infoMessages, 1)
	assert.Equal(t, "\033[31mtest message\033[0m", m.infoMessages[0])
}

func TestNewColored_Debug(t *testing.T) {
	m := &mockLogger{}
	c := NewColored(m, Green)

	require.NotNil(t, c)
	c.Debug("debug message")

	require.Len(t, m.debugMessages, 1)
	assert.Equal(t, "\033[32mdebug message\033[0m", m.debugMessages[0])
}

func TestNewColored_Warn(t *testing.T) {
	m := &mockLogger{}
	c := NewColored(m, Yellow)

	require.NotNil(t, c)
	c.Warn("warn message")

	require.Len(t, m.warnMessages, 1)
	assert.Equal(t, "\033[33mwarn message\033[0m", m.warnMessages[0])
}

func TestNewColored_Error(t *testing.T) {
	m := &mockLogger{}
	c := NewColored(m, Blue)

	require.NotNil(t, c)
	c.Error("error message")

	require.Len(t, m.errorMessages, 1)
	assert.Equal(t, "\033[34merror message\033[0m", m.errorMessages[0])
}

func TestNewColored_MultipleMessages(t *testing.T) {
	m := &mockLogger{}
	c := NewColored(m, Magenta)

	require.NotNil(t, c)
	c.Info("info 1")
	c.Info("info 2")
	c.Warn("warn 1")

	require.Len(t, m.infoMessages, 2)
	require.Len(t, m.warnMessages, 1)

	assert.Equal(t, "\033[35minfo 1\033[0m", m.infoMessages[0])
	assert.Equal(t, "\033[35minfo 2\033[0m", m.infoMessages[1])
	assert.Equal(t, "\033[35mwarn 1\033[0m", m.warnMessages[0])
}
