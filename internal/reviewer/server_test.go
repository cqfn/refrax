package reviewer

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

type mockReviewer struct {
	suggestions []domain.Suggestion
	err         error
}

func (r *mockReviewer) Review() (*domain.Artifacts, error) {
	if r.err != nil {
		return nil, r.err
	}
	return &domain.Artifacts{
		Descr:       &domain.Description{Text: "mock review"},
		Suggestions: r.suggestions,
	}, nil
}

func TestNewReviewer_CreatesValidInstance(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()

	reviewer := NewReviewer(ai, port, true, "echo tëst")

	require.NotNil(t, reviewer, "reviewer must not be nil")
}

func TestNewReviewer_SetsPort(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()

	reviewer := NewReviewer(ai, port, true, "echo ünïcödé")

	assert.Equal(t, port, reviewer.port, "reviewer port must match the provided value")
}

func TestReviewerStart_Success(t *testing.T) {
	ai := brain.NewMock()
	port, err := util.FreePort()
	require.NoError(t, err, "failed to get free port")
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	var listen error
	var shutdown error

	go func() { listen = reviewer.ListenAndServe() }()

	defer func() { shutdown = reviewer.Shutdown() }()
	require.NoError(t, shutdown, "unexpected error during shutdown")
	require.NoError(t, listen, "unexpected error during listen")
	_, ok := <-reviewer.Ready()
	assert.False(t, ok, "ready channel must be closed")
}

func TestReviewerStart_ServerStartError(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	reviewer.server = &mock{started: true}

	err := reviewer.ListenAndServe()

	require.Error(t, err, "expected error when server is already started")
}

func TestReviewerShutdown_SuccessAfterStart(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{started: true}
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	reviewer.server = server

	err := reviewer.Shutdown()

	require.NoError(t, err, "unexpected error during shutdown")
}

func TestReviewerShutdown_ServerNotStartedError(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{}
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	reviewer.server = server

	err := reviewer.Shutdown()

	require.Error(t, err, "expected error when server not started")
}

func TestReviewerThink_ReturnsMessageOnSuccess(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	reviewer.server = &mock{}
	reviewer.original = &mockReviewer{
		suggestions: []domain.Suggestion{*domain.NewSuggestion("add logging", "/path/Tëst.java")},
	}
	msg := protocol.NewMessage().WithMessageID("msg-ünïcödé").AddPart(protocol.NewText("review"))

	response, err := reviewer.think(context.Background(), msg)

	require.NoError(t, err, "unexpected error during think")
	require.NotNil(t, response, "response must not be nil")
}

func TestReviewerThink_HandlesContextCancellation(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	reviewer.server = &mock{}
	msg := protocol.NewMessage().WithMessageID("msg-cäncel").AddPart(protocol.NewText("review"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := reviewer.think(ctx, msg)

	require.Error(t, err, "expected error when context is canceled")
}

func TestReviewerThink_ReturnsErrorOnReviewFailure(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	reviewer.server = &mock{}
	reviewer.original = &mockReviewer{
		err: errors.New("review failed with ünïcödé"),
	}
	msg := protocol.NewMessage().WithMessageID("msg-ërror").AddPart(protocol.NewText("review"))

	_, err := reviewer.think(context.Background(), msg)

	require.Error(t, err, "expected error when review fails")
}

func TestAgentCard_ReturnsValidCard(t *testing.T) {
	port := randomPort()

	card := agentCard(port)

	require.NotNil(t, card, "agent card must not be nil")
}

func TestAgentCard_ContainsCorrectName(t *testing.T) {
	port := randomPort()

	card := agentCard(port)

	assert.Equal(t, "Reviewer Agent", card.Name, "agent card name must match")
}

func TestReviewerHandler_SetsHandler(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	reviewer := NewReviewer(ai, port, true, "echo tëst")
	called := false
	handler := func(next protocol.Handler, r *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
		called = true
		return nil, nil
	}

	reviewer.Handler(handler)

	assert.False(t, called, "handler must not be called during registration")
}

func TestLogSuggestions_LogsAllSuggestions(t *testing.T) {
	suggestions := []domain.Suggestion{
		*domain.NewSuggestion("add logging", "/path/Önë.java"),
		*domain.NewSuggestion("remove redundancy", "/path/Twö.java"),
	}
	logger := &mockLogger{messages: make([]string, 0)}

	logSuggestions(logger, suggestions)

	assert.Equal(t, 2, len(logger.messages), "must log all suggestions")
}

type mockLogger struct {
	messages []string
}

func (l *mockLogger) Info(msg string, args ...any) {
	l.messages = append(l.messages, msg)
}

func (l *mockLogger) Debug(msg string, args ...any) {
	l.messages = append(l.messages, msg)
}

func (l *mockLogger) Warn(msg string, args ...any) {
	l.messages = append(l.messages, msg)
}

func (l *mockLogger) Error(msg string, args ...any) {
	l.messages = append(l.messages, msg)
}
