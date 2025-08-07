// Package critic provides functionality for analyzing and critiquing Java code.
package critic

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/tool"
	"github.com/google/uuid"
)

// Critic represents the main struct responsible for analyzing code critiques.
type Critic struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
	port   int
	tools  []tool.Tool
}

const notFound = "No suggestions found"

const prompt = `Analyze the following Java code:

{{code}}

Identify possible improvements or flaws such as:

* grammar and spelling mistakes in comments,
* variables that can be inlined or removed without changing functionality,
* unnecessary comments inside methods,
* redundant code.

Don't suggest any changes that would alter the functionality of the code.
Don't suggest any changes that would require moving code parts between files (like extract class or extract an interface).
Don't suggest method renaming or class renaming.

Keep in mind the following imperfections with Java code, identified by automated static analysis system:

{{imperfections}}

Respond with a few most relevant and important suggestion for improvement, formatted as a few lines of text. 
If there are no suggestions or they are insignificant, respond with "{{not-found}}".
Do not include any explanations, summaries, or extra text.
`

// NewCritic creates and initializes a new instance of Critic.
func NewCritic(ai brain.Brain, port int, tools ...tool.Tool) *Critic {
	logger := log.NewPrefixed("critic", log.NewColored(log.Default(), log.Cyan))
	server := protocol.NewServer(agentCard(port), port)
	critic := &Critic{
		server: server,
		brain:  ai,
		log:    logger,
		port:   port,
		tools:  tools,
	}
	server.MsgHandler(critic.think)
	critic.log.Debug("preparing the Critic server on port %d with ai provider %s", port, ai)
	return critic
}

// ListenAndServe starts the Critic server and signals readiness via the provided channel.
func (c *Critic) ListenAndServe() error {
	c.log.Info("starting critic server on port %d...", c.port)
	var err error
	if err = c.server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
		return fmt.Errorf("failed to start critic server: %w", err)
	}
	return err
}

// Review sends the provided Java class to the Critic for analysis and returns suggested improvements.
func (c *Critic) Review(class domain.Class) ([]domain.Suggestion, error) {
	address := fmt.Sprintf("http://localhost:%d", c.port)
	c.log.Info("asking critic (%s) to lint the class...", address)
	critic := protocol.NewClient(address)
	msg := protocol.NewMessage().
		WithMessageID(uuid.NewString()).
		AddPart(protocol.NewText("lint class")).
		AddPart(protocol.NewFileBytes([]byte(class.Content())).WithMetadata("class-name", class.Name()).WithMetadata("class-path", class.Path()))
	resp, err := critic.SendMessage(protocol.NewMessageSendParams().WithMessage(msg))
	if err != nil {
		return nil, fmt.Errorf("failed to send message to critic: %w", err)
	}
	return domain.RespToSuggestions(resp), nil
}

// Shutdown gracefully shuts down the Critic server.
func (c *Critic) Shutdown() error {
	c.log.Info("stopping critic server...")
	if err := c.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop critic server: %w", err)
	}
	c.log.Info("critic server stopped successfully")
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
	c.log.Debug("received message: #%s", m.MessageID)
	tsk, err := domain.MsgToTask(m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task from message: %w", err)
	}
	class := tsk.Classes()[0]
	c.log.Info("received class %q for analysis", class.Name())
	replacer := strings.NewReplacer(
		"{{code}}", class.Content(),
		"{{imperfections}}", tool.NewCombined(c.tools...).Imperfections(),
		"{{not-found}}", notFound,
	)
	answer, err := c.brain.Ask(replacer.Replace(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from brain: %w", err)
	}
	res := protocol.NewMessage().WithMessageID(m.MessageID)
	suggestions := parseAnswer(answer)
	suggestions = associated(suggestions, class.Path())
	logSuggestions(c.log, suggestions)
	c.log.Info("found %d possible improvements", len(suggestions))
	for _, suggestion := range suggestions {
		c.log.Debug("suggestion: %s", suggestion)
		if strings.EqualFold(suggestion, notFound) {
			c.log.Info("no suggestions found")
		} else {
			res.AddPart(protocol.NewText(suggestion))
		}
	}
	c.log.Debug("sending response: %s", res.MessageID)
	return res, nil
}

func logSuggestions(logger log.Logger, suggestions []string) {
	for i, suggestion := range suggestions {
		logger.Info("suggestion #%d: %s", i+1, suggestion)
	}
}

func parseAnswer(answer string) []string {
	lines := strings.Split(strings.TrimSpace(answer), "\n")
	var suggestions []string
	for _, line := range lines {
		suggestion := strings.TrimSpace(line)
		if suggestion != "" {
			suggestions = append(suggestions, suggestion)
		}
	}
	return suggestions
}

func associated(suggestions []string, classPath string) []string {
	for i, suggestion := range suggestions {
		suggestions[i] = fmt.Sprintf("%s: %s", classPath, suggestion)
	}
	return suggestions
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.NewAgentCard().
		WithName("Critic Agent").
		WithDescription("Critic Description").
		WithURL(fmt.Sprintf("http://localhost:%d", port)).
		WithVersion("0.0.1").
		AddSkill("critic-java-code", "Critic Java Code", "Give a reasonable critique on Java code")
}
