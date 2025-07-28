package critic

import (
	"context"
	"errors"
	"testing"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mock struct {
	started bool
	closed  bool
	handler protocol.MsgHandler
	ready   chan bool
}

func (m *mock) ListenAndServe() error {
	if m.started {
		return errors.New("server already started")
	}
	reeady := make(chan bool, 1)
	reeady <- true
	close(reeady)
	m.ready = reeady
	m.started = true
	return nil
}

func (m *mock) Shutdown() error {
	if !m.started {
		return errors.New("server not started")
	}
	if m.closed {
		return errors.New("server already closed")
	}
	m.closed = true
	return nil
}

func (m *mock) Handler(_ protocol.Handler) {
}

func (m *mock) MsgHandler(handler protocol.MsgHandler) {
	m.handler = handler
}

func (m *mock) Ready() <-chan bool {
	return m.ready
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
	var listen error
	var shutdown error

	go func() { listen = critic.ListenAndServe() }()

	defer func() { shutdown = critic.Shutdown() }()
	require.NoError(t, shutdown)
	require.NoError(t, listen)
	_, ok := <-critic.Ready()
	assert.False(t, ok)
}

func TestCriticStart_ServerStartError(t *testing.T) {
	ai := brain.NewMock()
	critic := NewCritic(ai, 18081)
	critic.server = &mock{started: true}

	err := critic.ListenAndServe()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start critic server")
}

func TestCriticClose_Success(t *testing.T) {
	ai := brain.NewMock()
	server := &mock{started: true}
	critic := NewCritic(ai, 18081)
	critic.server = server

	err := critic.Shutdown()

	require.NoError(t, err)
	assert.True(t, server.closed)
}

func TestCriticClose_ServerNotStartedError(t *testing.T) {
	ai := brain.NewMock()
	server := &mock{}
	critic := NewCritic(ai, 18081, NewMockToolEmpty())
	critic.server = server

	err := critic.Shutdown()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stop critic server")
}

func TestCriticThink_ReturnsMessage(t *testing.T) {
	ai := brain.NewMock()
	server := &mock{}
	critic := NewCritic(ai, 18081, NewMockToolEmpty())
	critic.server = server
	msg := protocol.NewMessageBuilder().
		MessageID("msg-123").
		Build()

	response, err := critic.think(context.Background(), msg)

	require.NoError(t, err)
	assert.Equal(t, msg.MessageID, response.MessageID)
}
