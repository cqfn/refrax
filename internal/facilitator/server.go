package facilitator

import (
	"fmt"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Facilitator struct {
	server protocol.Server
}

func NewFacilitator(ai string, port int) (*Facilitator, error) {
	log.Debug("preparing server on port %d with ai provider %s", port, ai)
	server, error := protocol.NewCustomServer(agentCard(port), protocol.LogRequest, port)
	if error != nil {
		return nil, fmt.Errorf("failed to create A2A server: %w", error)
	}
	return &Facilitator{
		server: server,
	}, nil
}

func StartServer(ai string, port int) error {
	log.Debug("starting refrax server on port %d with ai provider %s", port, ai)
	facilitator, err := NewFacilitator(ai, port)
	if err != nil {
		return fmt.Errorf("failed to create facilitator: %w", err)
	}
	return facilitator.Start(make(chan struct{}))
}

func (f *Facilitator) Start(ready chan<- struct{}) error {
	log.Debug("starting facilitator server")
	if err := f.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
	log.Info("facilitator server started successfully")
	return nil
}

func (f *Facilitator) Close() error {
	log.Debug("stopping facilitator server")
	if err := f.server.Close(); err != nil {
		return fmt.Errorf("failed to stop facilitator server: %w", err)
	}
	log.Info("facilitator server stopped successfully")
	return nil
}

func agentCard(port int) protocol.AgentCard {
	return protocol.Card().
		Name("Facilitator Agent").
		Description("An agent that facilitates talk between critic and fixer").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("facilitate-discussion", "Refactor Java Projects", "Facilitate discussion on code refactoring").
		Build()
}
