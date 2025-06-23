package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrpcServer_Creates_Successfully(t *testing.T) {
	port := 16745
	server, err := NewTrpcServer(mockCard(port), port)

	require.NoError(t, err, "should create a new TRPC server without error")
	assert.NotNil(t, server, "server should not be nil")
}
