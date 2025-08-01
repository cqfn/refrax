// Package facilitator provides functionality for facilitating interactions between
// critic and fixer agents in a code refactoring process.
package facilitator

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

// A2AFacilitator facilitates communication between the critic and fixer agents.
type A2AFacilitator struct {
	server   protocol.Server
	log      log.Logger
	port     int
	original domain.Facilitator
}

// NewFacilitator creates a new instance of Facilitator to manage communication between agents.
func NewFacilitator(ai brain.Brain, critic domain.Critic, fixer domain.Fixer, port int) *A2AFacilitator {
	logger := log.NewPrefixed("facilitator", log.NewColored(log.Default(), log.Yellow))
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewServer(agentCard(port), port)
	facilitator := &A2AFacilitator{
		server: server,
		log:    logger,
		port:   port,
		original: &agent{
			brain:  ai,
			log:    logger,
			critic: critic,
			fixer:  fixer,
		},
	}
	server.MsgHandler(facilitator.think)
	return facilitator
}

// Refactor sends a refactoring request to the facilitator server and returns the refactored classes.
func (f *A2AFacilitator) Refactor(task domain.Task) ([]domain.Class, error) {
	client := protocol.NewClient(fmt.Sprintf("http://localhost:%d", f.port))
	resp, err := client.SendMessage(
		protocol.NewMessageSendParams().WithMessage(domain.TaskToMsg(task)))
	if err != nil {
		return nil, fmt.Errorf("failed to send refactoring request: %w", err)
	}
	return domain.RespToClasses(resp)
}

// ListenAndServe starts the facilitator server and prepares it for handling requests.
func (f *A2AFacilitator) ListenAndServe() error {
	f.log.Info("starting facilitator server on port %d...", f.port)
	var err error
	if err = f.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
	return err
}

// Shutdown stops the facilitator server and releases resources.
func (f *A2AFacilitator) Shutdown() error {
	f.log.Info("stopping facilitator server...")
	if err := f.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop facilitator server: %w", err)
	}
	f.log.Info("facilitator server stopped successfully")
	return nil
}

// Ready returns a channel that signals when the facilitator server is ready to accept requests.
func (f *A2AFacilitator) Ready() <-chan bool {
	return f.server.Ready()
}

// Handler sets the message handler for the facilitator server.
func (f *A2AFacilitator) Handler(handler protocol.Handler) {
	f.server.Handler(handler)
}

func (f *A2AFacilitator) think(ctx context.Context, m *protocol.Message) (*protocol.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return f.thinkLong(m)
	}
}

func (f *A2AFacilitator) thinkLong(m *protocol.Message) (*protocol.Message, error) {
	task, err := domain.MsgToTask(m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	resp, err := f.original.Refactor(task)
	if err != nil {
		return nil, fmt.Errorf("failed to refactor task: %w", err)
	}
	log.Debug("number of processed classes: %d", len(resp))
	return domain.ClassesToMsg(resp), nil
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.NewAgentCard().
		WithName("Facilitator Agent").
		WithDescription("An agent that facilitates talk between critic and fixer").
		WithURL(fmt.Sprintf("http://localhost:%d", port)).
		WithVersion("0.0.1").
		AddSkill("facilitate-discussion", "Refactor Java Projects", "Facilitate discussion on code refactoring")
}
