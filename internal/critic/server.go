package critic

import (
	"fmt"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Critic struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
}

func NewCritic(ai brain.Brain, port int) (*Critic, error) {
	logger := log.NewPrefixed("critic", log.Default())
	server := protocol.NewCustomServer(agentCard(port), port)
	critic := &Critic{
		server: server,
		brain:  ai,
		log:    logger,
	}
	server.SetHandler(critic.think)
	critic.log.Debug("preparing the Critic server on port %d with ai provider %s", port, ai)
	return critic, nil
}

func (c *Critic) Start(ready chan<- struct{}) error {
	c.log.Debug("starting critic server")
	if err := c.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start critic server: %w", err)
	}
	c.log.Info("critic server started successfully")
	return nil
}

func (c *Critic) Close() error {
	c.log.Debug("stopping critic server")
	if err := c.server.Close(); err != nil {
		return fmt.Errorf("failed to stop critic server: %w", err)
	}
	c.log.Info("critic server stopped successfully")
	return nil
}

func (c *Critic) think(m *protocol.Message) (*protocol.Message, error) {
	c.log.Info("received message id:  %s", m.MessageID)
	c.log.Debug("message parts: %v", m.Parts)
	// Here you would implement the logic to critique the Java code
	// For now, we just return the message back
	// use brain
	return m, nil
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
