package reviewer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

// A2AReviewer represents a reviewer agent responsible for reviewing changes.
type A2AReviewer struct {
	server   protocol.Server
	log      log.Logger
	port     int
	original domain.Reviewer
}

// NewReviewer creates a new instance of A2AReviewer.
func NewReviewer(ai brain.Brain, port int, cmds ...string) *A2AReviewer {
	logger := log.NewPrefixed("reviewer", log.NewColored(log.Default(), log.Orange))
	logger.Debug("preparing server on port %d", port)
	server := protocol.NewServer(agentCard(port), port)
	reviewer := &A2AReviewer{
		server: server,
		log:    logger,
		port:   port,
		original: &agent{
			logger: logger,
			cmds:   cmds,
			ai:     ai,
		},
	}
	server.MsgHandler(reviewer.think)
	return reviewer
}

// Review sends a request for review and returns suggestions.
func (r *A2AReviewer) Review() (*domain.Artifacts, error) {
	client := protocol.NewClient(fmt.Sprintf("http://localhost:%d", r.port))
	resp, err := client.SendMessage(
		protocol.NewMessageSendParams().WithMessage(protocol.NewMessage().AddPart(protocol.NewText("review"))),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send review request: %w", err)
	}
	return domain.UnmarshalArtifacts(resp.Result.(*protocol.Message))
}

// ListenAndServe starts the reviewer server and listens for incoming requests.
func (r *A2AReviewer) ListenAndServe() error {
	r.log.Info("starting reviewer server on port %d...", r.port)
	var err error
	if err = r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start reviewer server: %w", err)
	}
	return err
}

// Shutdown gracefully shuts down the reviewer server.
func (r *A2AReviewer) Shutdown() error {
	r.log.Info("stopping reviewer server...")
	if err := r.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop reviewer server: %w", err)
	}
	r.log.Info("reviewer server stopped successfully")
	return nil
}

// Ready returns a channel that signals when the server is ready to handle requests.
func (r *A2AReviewer) Ready() <-chan bool {
	return r.server.Ready()
}

// Handler sets a custom handler for the reviewer server.
func (r *A2AReviewer) Handler(handler protocol.Handler) {
	r.server.Handler(handler)
}

func (r *A2AReviewer) think(ctx context.Context, m *protocol.Message) (*protocol.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return r.thinkLong(m)
	}
}

func (r *A2AReviewer) thinkLong(_ *protocol.Message) (*protocol.Message, error) {
	artifacts, err := r.original.Review()
	if err != nil {
		return nil, fmt.Errorf("failed to  task: %w", err)
	}
	suggestions := artifacts.Suggestions
	r.log.Info("number of processed suggestions: %d", len(suggestions))
	logSuggestions(r.log, suggestions)
	r.log.Info("found %d possible fixes", len(suggestions))
	return artifacts.Marshal().Message, nil
}

func logSuggestions(logger log.Logger, suggestions []domain.Suggestion) {
	for i, suggestion := range suggestions {
		logger.Info("suggestion #%d: %s", i+1, suggestion)
	}
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.NewAgentCard().
		WithName("Reviewer Agent").
		WithDescription("An agent that checks whether the project is stable and changes made haven't break anything").
		WithURL(fmt.Sprintf("http://localhost:%d", port)).
		WithVersion("0.0.1").
		AddSkill("review-changes", "Review Changes", "Review changes made to be sure they don't break the project")
}
