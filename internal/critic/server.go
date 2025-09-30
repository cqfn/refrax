package critic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/tool"
)

// Critic represents the main struct responsible for analyzing code critiques.
type Critic struct {
	server protocol.Server
	log    log.Logger
	port   int
	agent  *agent
}

// NewCritic creates and initializes a new instance of Critic.
func NewCritic(ai brain.Brain, port int, colorless bool, tools ...tool.Tool) *Critic {
	logger := log.New("critic", log.Cyan, colorless)
	server := protocol.NewServer(agentCard(port), port)
	critic := &Critic{
		server: server,
		log:    logger,
		port:   port,
		agent:  &agent{brain: ai, log: logger, tools: tools},
	}
	server.MsgHandler(critic.think)
	critic.log.Debug("Preparing the Critic server on port %d with ai provider %s", port, ai)
	return critic
}

// ListenAndServe starts the Critic server and signals readiness via the provided channel.
func (c *Critic) ListenAndServe() error {
	c.log.Info("Starting critic server on port %d...", c.port)
	var err error
	if err = c.server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
		return fmt.Errorf("failed to start critic server: %w", err)
	}
	return err
}

// Review sends the provided Java class to the Critic for analysis and returns suggested improvements.
func (c *Critic) Review(job *domain.Job) (*domain.Artifacts, error) {
	address := fmt.Sprintf("http://localhost:%d", c.port)
	c.log.Debug("Asking critic (%s) to lint the class...", address)
	critic := protocol.NewClient(address)
	resp, err := critic.SendMessage(job.Marshal())
	if err != nil {
		return nil, fmt.Errorf("failed to send message to critic: %w", err)
	}
	return domain.UnmarshalArtifacts(resp.Result.(*protocol.Message))
}

// Shutdown gracefully shuts down the Critic server.
func (c *Critic) Shutdown() error {
	c.log.Info("Stopping critic server...")
	if err := c.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop critic server: %w", err)
	}
	c.log.Info("Critic server stopped successfully")
	return nil
}

// Handler sets the message handler for the Critic server.
func (c *Critic) Handler(handler protocol.Handler) {
	c.server.Handler(handler)
}

// Ready returns a channel that signals when the Critic server is ready to accept requests.
func (c *Critic) Ready() <-chan bool {
	return c.server.Ready()
}

func (c *Critic) think(ctx context.Context, m *protocol.Message) (*protocol.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return c.thinkLong(m)
	}
}

func (c *Critic) thinkLong(m *protocol.Message) (*protocol.Message, error) {
	c.log.Debug("Received message: #%s", m.MessageID)
	tsk, err := domain.UnmarshalJob(m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task from message: %w", err)
	}
	artifacts, err := c.agent.Review(tsk)
	if err != nil {
		return nil, fmt.Errorf("failed to review the task: %w", err)
	}
	return artifacts.Marshal().Message, err
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.NewAgentCard().
		WithName("Critic Agent").
		WithDescription("Critic Description").
		WithURL(fmt.Sprintf("http://localhost:%d", port)).
		WithVersion("0.0.1").
		AddSkill("critic-java-code", "Critic Java Code", "Give a reasonable critique on Java code")
}
