package critic

import (
	"errors"
	"testing"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockServer struct {
	started bool
	closed  bool
	handler protocol.MsgHandler
}

func (m *MockServer) Start(ready chan<- struct{}) error {
	if m.started {
		return errors.New("server already started")
	}
	m.started = true
	close(ready)
	return nil
}

func (m *MockServer) Close() error {
	if !m.started {
		return errors.New("server not started")
	}
	if m.closed {
		return errors.New("server already closed")
	}
	m.closed = true
	return nil
}

func (m *MockServer) Handler(_ protocol.Handler) {
}

func (m *MockServer) MsgHandler(handler protocol.MsgHandler) {
	m.handler = handler
}

type MockBrain struct{}

func TestNewCritic_Success(t *testing.T) {
	ai := brain.NewMock()

	critic := NewCritic(ai, 18081)
	require.NotNil(t, critic)
	assert.Equal(t, ai, critic.brain)
}

func TestCriticStart_Success(t *testing.T) {
	ai := brain.NewMock()
	port, err := protocol.FreePort()
	require.NoError(t, err)
	critic := NewCritic(ai, port)
	ready := make(chan struct{})

	go func() { err = critic.Start(ready) }()

	defer func() { err = critic.Close() }()
	require.NoError(t, err)
	_, ok := <-ready
	assert.False(t, ok)
}

func TestCriticStart_ServerStartError(t *testing.T) {
	ai := brain.NewMock()
	critic := NewCritic(ai, 18081)
	critic.server = &MockServer{started: true}
	ready := make(chan struct{})

	err := critic.Start(ready)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start critic server")
}

func TestCriticClose_Success(t *testing.T) {
	ai := brain.NewMock()
	server := &MockServer{started: true}
	critic := NewCritic(ai, 18081)
	critic.server = server

	err := critic.Close()

	require.NoError(t, err)
	assert.True(t, server.closed)
}

func TestCriticClose_ServerNotStartedError(t *testing.T) {
	ai := brain.NewMock()
	server := &MockServer{}
	critic := NewCritic(ai, 18081, NewMockToolEmpty())
	critic.server = server

	err := critic.Close()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stop critic server")
}

func TestCriticThink_ReturnsMessage(t *testing.T) {
	ai := brain.NewMock()
	server := &MockServer{}
	critic := NewCritic(ai, 18081, NewMockToolEmpty())
	critic.server = server
	msg := protocol.NewMessageBuilder().
		MessageID("msg-123").
		Build()

	response, err := critic.think(msg)

	require.NoError(t, err)
	assert.Equal(t, msg.MessageID, response.MessageID)
}
