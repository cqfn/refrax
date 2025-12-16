package fixer

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"testing"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomPort() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return int(n.Int64()) + 20000
}

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

func TestNewFixer_CreatesValidInstance(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()

	fixer := NewFixer(ai, port, false)

	require.NotNil(t, fixer, "fixer must not be nil")
}

func TestNewFixer_SetsPort(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()

	fixer := NewFixer(ai, port, true)

	assert.Equal(t, port, fixer.port, "fixer port must match the provided value")
}

func TestNewFixer_SetsBrain(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()

	fixer := NewFixer(ai, port, false)

	assert.Equal(t, ai, fixer.brain, "fixer brain must match the provided AI")
}

func TestFixerStart_Success(t *testing.T) {
	ai := brain.NewMock()
	port, err := util.FreePort()
	require.NoError(t, err, "failed to get free port")
	fixer := NewFixer(ai, port, true)
	var listen error
	var shutdown error

	go func() { listen = fixer.ListenAndServe() }()

	defer func() { shutdown = fixer.Shutdown() }()
	require.NoError(t, shutdown, "unexpected error during shutdown")
	require.NoError(t, listen, "unexpected error during listen")
	_, ok := <-fixer.Ready()
	assert.False(t, ok, "ready channel must be closed")
}

func TestFixerStart_ServerStartError(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	fixer := NewFixer(ai, port, true)
	fixer.server = &mock{started: true}

	err := fixer.ListenAndServe()

	require.Error(t, err, "expected an error when server is already started")
}

func TestFixerShutdown_SuccessAfterStart(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{started: true}
	fixer := NewFixer(ai, port, true)
	fixer.server = server

	err := fixer.Shutdown()

	require.NoError(t, err, "unexpected error during shutdown")
}

func TestFixerShutdown_ServerNotStartedError(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{}
	fixer := NewFixer(ai, port, true)
	fixer.server = server

	err := fixer.Shutdown()

	require.Error(t, err, "expected error when server not started")
}

func TestFixerThink_ReturnsMessageOnSuccess(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{}
	fixer := NewFixer(ai, port, true)
	fixer.server = server
	class := domain.NewInMemoryClass("Tëst.java", "/path/Tëst.java", "public class Tëst {}")
	job := &domain.Job{
		Descr:       &domain.Description{Text: "fix the class"},
		Classes:     []domain.Class{class},
		Suggestions: []domain.Suggestion{*domain.NewSuggestion("add logging", "/path/Tëst.java")},
	}
	msg := job.Marshal().Message

	response, err := fixer.think(context.Background(), msg)

	require.NoError(t, err, "unexpected error during think")
	require.NotNil(t, response, "response must not be nil")
}

func TestFixerThink_HandlesContextCancellation(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	fixer := NewFixer(ai, port, true)
	fixer.server = &mock{}
	class := domain.NewInMemoryClass("Fïle.java", "/path/Fïle.java", "class Fïle {}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "fix the class"},
		Classes: []domain.Class{class},
	}
	msg := job.Marshal().Message
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := fixer.think(ctx, msg)

	require.Error(t, err, "expected error when context is canceled")
}

func TestClean_RemovesJavaCodeFences(t *testing.T) {
	input := "```java\npublic class Tëst {}\n```"

	result := clean(input)

	assert.Equal(t, "\npublic class Tëst {}\n", result, "code fences must be removed")
}

func TestClean_RemovesPlainCodeFences(t *testing.T) {
	input := "```\nclass Ünïcödé {}\n```"

	result := clean(input)

	assert.Equal(t, "\nclass Ünïcödé {}\n", result, "plain code fences must be removed")
}

func TestClean_PreservesContentWithoutFences(t *testing.T) {
	input := "public class NöFence {}"

	result := clean(input)

	assert.Equal(t, "public class NöFence {}", result, "content without fences must be preserved")
}

func TestAgentCard_ReturnsValidCard(t *testing.T) {
	port := randomPort()

	card := agentCard(port)

	require.NotNil(t, card, "agent card must not be nil")
}

func TestAgentCard_ContainsCorrectName(t *testing.T) {
	port := randomPort()

	card := agentCard(port)

	assert.Equal(t, "Fixer Agent", card.Name, "agent card name must match")
}

func TestFixerHandler_SetsHandler(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	fixer := NewFixer(ai, port, true)
	called := false
	handler := func(next protocol.Handler, r *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
		called = true
		return nil, nil
	}

	fixer.Handler(handler)

	assert.False(t, called, "handler must not be called during registration")
}

func TestThinkChan_ReturnsChannelWithResult(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	fixer := NewFixer(ai, port, true)
	class := domain.NewInMemoryClass("Chännel.java", "/Chännel.java", "class Chännel{}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "fix the class"},
		Classes: []domain.Class{class},
	}
	msg := job.Marshal().Message

	ch := fixer.thinkChan(msg)
	result := <-ch

	require.NoError(t, result.err, "unexpected error from thinkChan")
}
