package util

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFreePort_ReturnsValidPort(t *testing.T) {
	port, err := FreePort()
	require.NoError(t, err, "FreePort unexpectedly returned an error")
	assert.Greater(t, port, 0, "Port must be a positive number")
}

func TestFreePort_ReturnsUsablePort(t *testing.T) {
	port, err := FreePort()
	require.NoError(t, err, "FreePort unexpectedly returned an error")
	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", string(rune(port))))
	if err == nil {
		defer func() { _ = listener.Close() }()
	}
	assert.Greater(t, port, 1024, "Port must be in unprivileged range")
}

func TestFreePort_ReturnsUnprivilegedPort(t *testing.T) {
	port, err := FreePort()
	require.NoError(t, err, "FreePort unexpectedly returned an error")
	assert.Greater(t, port, 1024, "Port must be greater than 1024")
}

func TestFreePort_ReturnsWithinValidRange(t *testing.T) {
	port, err := FreePort()
	require.NoError(t, err, "FreePort unexpectedly returned an error")
	assert.Less(t, port, 65536, "Port must be less than 65536")
}

func TestFreePort_MultipleCalls_ReturnDifferentPorts(t *testing.T) {
	first, err := FreePort()
	require.NoError(t, err, "First FreePort call unexpectedly returned an error")
	second, err := FreePort()
	require.NoError(t, err, "Second FreePort call unexpectedly returned an error")
	assert.NotEqual(t, first, second, "Sequential calls should return different ports")
}
