// Package facilitator provides functionality for facilitating interactions between
// critic and fixer agents in a code refactoring process.
package facilitator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

// Facilitator facilitates communication between the critic and fixer agents.
type Facilitator struct {
	server   protocol.Server
	log      log.Logger
	port     int
	original agent
}

// NewFacilitator creates a new instance of Facilitator to manage communication between agents.
func NewFacilitator(ai brain.Brain, port, criticPort, fixerPort int) *Facilitator {
	logger := log.NewPrefixed("facilitator", log.NewColored(log.Default(), log.Yellow))
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewServer(agentCard(port), port)
	facilitator := &Facilitator{
		server: server,
		log:    logger,
		port:   port,
		original: agent{
			brain:      ai,
			log:        logger,
			criticPort: criticPort,
			fixerPort:  fixerPort,
		},
	}
	server.MsgHandler(facilitator.think)
	return facilitator
}

// ListenAndServe starts the facilitator server and prepares it for handling requests.
func (f *Facilitator) ListenAndServe() error {
	f.log.Info("starting facilitator server on port %d...", f.port)
	var err error
	if err = f.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
	return err
}

// Shutdown stops the facilitator server and releases resources.
func (f *Facilitator) Shutdown() error {
	f.log.Info("stopping facilitator server...")
	if err := f.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop facilitator server: %w", err)
	}
	f.log.Info("facilitator server stopped successfully")
	return nil
}

// Ready returns a channel that signals when the facilitator server is ready to accept requests.
func (f *Facilitator) Ready() <-chan bool {
	return f.server.Ready()
}

// Handler sets the message handler for the facilitator server.
func (f *Facilitator) Handler(handler protocol.Handler) {
	f.server.Handler(handler)
}

func (f *Facilitator) think(ctx context.Context, m *protocol.Message) (*protocol.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return f.thinkLong(m)
	}
}

func (f *Facilitator) thinkLong(m *protocol.Message) (*protocol.Message, error) {
	if f.isRefactoringTask(m) {
		return f.refactor(m)
	}
	f.log.Warn("received a message that is not related to refactoring. ignoring.")
	return nil, fmt.Errorf("received a message that is not related to refactoring")
}

func (f *Facilitator) isRefactoringTask(m *protocol.Message) bool {
	res := false
	if len(m.Parts) > 0 {
		for _, part := range m.Parts {
			if part.PartKind() == protocol.PartKindText {
				task := part.(*protocol.TextPart).Text
				if strings.Contains(task, "refactor the project") {
					res = true
					break
				}
			}
		}
	}
	if !res {
		f.log.Warn("received a task that is not related to refactoring. ignoring.")
	}
	return res
}

func (f *Facilitator) refactor(m *protocol.Message) (*protocol.Message, error) {
	task, err := ParseTask(m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	resp, err := f.original.refactor(task)
	if err != nil {
		return nil, fmt.Errorf("failed to refactor task: %w", err)
	}
	log.Debug("number of processed classes: %d", len(resp))
	log.Debug("refactored class: %s", resp[0].Name())
	msg := CompileMessage(resp)
	log.Debug("sending response: %s", msg)
	return msg, nil
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.Card().
		Name("Facilitator Agent").
		Description("An agent that facilitates talk between critic and fixer").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("facilitate-discussion", "Refactor Java Projects", "Facilitate discussion on code refactoring").
		Build()
}
