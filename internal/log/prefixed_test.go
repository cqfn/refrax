package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrefixed(t *testing.T) {
	mock := NewMock()
	prefix := "prefix"
	pref := NewPrefixed(prefix, mock)
	require.NotNil(t, pref)

	typed, ok := pref.(*prefixed)
	require.True(t, ok)
	assert.Equal(t, prefix, typed.prefix)
	assert.Equal(t, mock, typed.original)
}

func TestPrefixed_Info(t *testing.T) {
	mock := NewMock().(*Mock)
	prefix := "test-prefix"
	pref := NewPrefixed(prefix, mock)

	pref.Info("info message")

	require.Len(t, mock.Messages, 1)
	assert.Equal(t, "mock info: test-prefix:      info message", mock.Messages[0])
}

func TestPrefixed_Debug(t *testing.T) {
	mock := NewMock().(*Mock)
	prefix := "debug-prefix"
	pref := NewPrefixed(prefix, mock)

	pref.Debug("debug message")

	require.Len(t, mock.Messages, 1)
	assert.Equal(t, "mock debug: debug-prefix:     debug message", mock.Messages[0])
}

func TestPrefixed_Warn(t *testing.T) {
	mock := NewMock().(*Mock)
	prefix := "warn-prefix"
	pref := NewPrefixed(prefix, mock)

	pref.Warn("warn message")

	require.Len(t, mock.Messages, 1)
	assert.Equal(t, "mock warn: warn-prefix:      warn message", mock.Messages[0])
}

func TestPrefixed_Error(t *testing.T) {
	mock := NewMock().(*Mock)
	prefix := "error-prefix"
	pref := NewPrefixed(prefix, mock)

	pref.Error("error message")

	require.Len(t, mock.Messages, 1)
	assert.Equal(t, "mock error: error-prefix:     error message", mock.Messages[0])
}

func TestPrefixed_WithArgs(t *testing.T) {
	mock := NewMock().(*Mock)
	prefix := "arg-prefix"
	pref := NewPrefixed(prefix, mock)

	pref.Info("message with args: %s, %d", "arg1", 42)
	pref.Debug("debug with args: %s, %t", "arg2", true)
	pref.Warn("warn with args: %s", "arg3")
	pref.Error("error with args: %s %.2f", "arg4", 3.14)

	// Assert each message was logged with the correct prefix
	require.Len(t, mock.Messages, 4)
	assert.Equal(t, "mock info: arg-prefix:       message with args: arg1, 42", mock.Messages[0])
	assert.Equal(t, "mock debug: arg-prefix:       debug with args: arg2, true", mock.Messages[1])
	assert.Equal(t, "mock warn: arg-prefix:       warn with args: arg3", mock.Messages[2])
	assert.Equal(t, "mock error: arg-prefix:       error with args: arg4 3.14", mock.Messages[3])
}
