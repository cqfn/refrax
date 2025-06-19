package critic

import (
	"fmt"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Critic struct {
	server protocol.Server
}

func NewCritic(ai string, port int) (*Critic, error) {
	log.Debug("preparing server on port %d with ai provider %s", port, ai)
	server, error := protocol.NewCustomServer(agentCard(port), protocol.LogRequest, port)
	if error != nil {
		return nil, fmt.Errorf("failed to create A2A server: %w", error)
	}
	return &Critic{
		server: server,
	}, nil
}

func (f *Critic) Start(ready chan<- struct{}) error {
	log.Debug("starting critic server")
	if err := f.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start critic server: %w", err)
	}
	log.Info("critic server started successfully")
	return nil
}

func (f *Critic) Close() error {
	log.Debug("stopping critic server")
	if err := f.server.Close(); err != nil {
		return fmt.Errorf("failed to stop critic server: %w", err)
	}
	log.Info("critic server stopped successfully")
	return nil
}

func agentCard(port int) protocol.AgentCard {
	return protocol.Card().
		Name("Critic Agent").
		Description("Critic Description").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("critic-java-code", "Critic Java Code", "Give a reasonable critique on Java code").
		Build()
}
