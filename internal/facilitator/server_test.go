package facilitator

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

type mockCritic struct {
	suggestions []domain.Suggestion
	err         error
}

func (c *mockCritic) Review(_ *domain.Job) (*domain.Artifacts, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &domain.Artifacts{
		Descr:       &domain.Description{Text: "mock review"},
		Suggestions: c.suggestions,
	}, nil
}

type mockFixer struct {
	classes []domain.Class
	err     error
}

func (f *mockFixer) Fix(job *domain.Job) (*domain.Artifacts, error) {
	if f.err != nil {
		return nil, f.err
	}
	if len(f.classes) > 0 {
		return &domain.Artifacts{
			Descr:   &domain.Description{Text: "mock fix"},
			Classes: f.classes,
		}, nil
	}
	return &domain.Artifacts{
		Descr:   &domain.Description{Text: "mock fix"},
		Classes: job.Classes,
	}, nil
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
		Descr:       &domain.Description{Text: "mock reviewer result"},
		Suggestions: r.suggestions,
	}, nil
}

type mockFacilitator struct {
	classes []domain.Class
	err     error
}

func (f *mockFacilitator) Refactor(_ *domain.Job) (*domain.Artifacts, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &domain.Artifacts{
		Descr:   &domain.Description{Text: "mock refactor"},
		Classes: f.classes,
	}, nil
}

func TestNewFacilitator_CreatesValidInstance(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}

	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)

	require.NotNil(t, facilitator, "facilitator must not be nil")
}

func TestNewFacilitator_SetsPort(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}

	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)

	assert.Equal(t, port, facilitator.port, "facilitator port must match the provided value")
}

func TestFacilitatorStart_Success(t *testing.T) {
	ai := brain.NewMock()
	port, err := util.FreePort()
	require.NoError(t, err, "failed to get free port")
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	var listen error
	var shutdown error

	go func() { listen = facilitator.ListenAndServe() }()

	defer func() { shutdown = facilitator.Shutdown() }()
	require.NoError(t, shutdown, "unexpected error during shutdown")
	require.NoError(t, listen, "unexpected error during listen")
	_, ok := <-facilitator.Ready()
	assert.False(t, ok, "ready channel must be closed")
}

func TestFacilitatorStart_ServerStartError(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	facilitator.server = &mock{started: true}

	err := facilitator.ListenAndServe()

	require.Error(t, err, "expected error when server is already started")
}

func TestFacilitatorShutdown_SuccessAfterStart(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{started: true}
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	facilitator.server = server

	err := facilitator.Shutdown()

	require.NoError(t, err, "unexpected error during shutdown")
}

func TestFacilitatorShutdown_ServerNotStartedError(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	server := &mock{}
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	facilitator.server = server

	err := facilitator.Shutdown()

	require.Error(t, err, "expected error when server not started")
}

func TestFacilitatorThink_ReturnsMessageOnSuccess(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 0)
	facilitator.server = &mock{}
	facilitator.original = &mockFacilitator{
		classes: []domain.Class{domain.NewInMemoryClass("Tëst.java", "/path/Tëst.java", "class Tëst{}")},
	}
	class := domain.NewInMemoryClass("Tëst.java", "/path/Tëst.java", "public class Tëst {}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "refactor the project"},
		Classes: []domain.Class{class},
	}
	msg := job.Marshal().Message

	response, err := facilitator.think(context.Background(), msg)

	require.NoError(t, err, "unexpected error during think")
	require.NotNil(t, response, "response must not be nil")
}

func TestFacilitatorThink_HandlesContextCancellation(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	facilitator.server = &mock{}
	class := domain.NewInMemoryClass("Fïle.java", "/path/Fïle.java", "class Fïle {}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "refactor the project"},
		Classes: []domain.Class{class},
	}
	msg := job.Marshal().Message
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := facilitator.think(ctx, msg)

	require.Error(t, err, "expected error when context is canceled")
}

func TestFacilitatorThink_ReturnsErrorOnRefactorFailure(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	facilitator.server = &mock{}
	facilitator.original = &mockFacilitator{
		err: errors.New("refactor failed with ünïcödé"),
	}
	class := domain.NewInMemoryClass("Ërror.java", "/path/Ërror.java", "class Ërror{}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "refactor the project"},
		Classes: []domain.Class{class},
	}
	msg := job.Marshal().Message

	_, err := facilitator.think(context.Background(), msg)

	require.Error(t, err, "expected error when refactor fails")
}

func TestAgentCard_ReturnsValidCard(t *testing.T) {
	port := randomPort()

	card := agentCard(port)

	require.NotNil(t, card, "agent card must not be nil")
}

func TestAgentCard_ContainsCorrectName(t *testing.T) {
	port := randomPort()

	card := agentCard(port)

	assert.Equal(t, "Facilitator Agent", card.Name, "agent card name must match")
}

func TestFacilitatorHandler_SetsHandler(t *testing.T) {
	ai := brain.NewMock()
	port := randomPort()
	critic := &mockCritic{}
	fixer := &mockFixer{}
	reviewer := &mockReviewer{}
	facilitator := NewFacilitator(ai, critic, fixer, reviewer, port, true, 1)
	called := false
	handler := func(next protocol.Handler, r *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
		called = true
		return nil, nil
	}

	facilitator.Handler(handler)

	assert.False(t, called, "handler must not be called during registration")
}
