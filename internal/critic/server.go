// Package critic provides functionality for analyzing and critiquing Java code.
package critic

import (
	"context"
	"encoding/base64"
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

const prompt = `Analyze the following Java code:

{{code}}

Identify possible improvements or flaws such as:

* grammar and spelling mistakes in comments,
* variables that can be inlined or removed without changing functionality,
* unnecessary comments inside methods,
* redundant code,
* non-idiomatic patterns.

Keep in mind the following imperfections with Java code, identified by automated static analysis system:

{{imperfections}}

Respond with a single most relevant and important suggestion for improvement, formatted as a single line of text. 
If there are no suggestions or they are insignificant, respond with "No suggestions found".
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
	msg := protocol.NewMessageBuilder().
		MessageID(uuid.NewString()).
		Part(protocol.NewText("lint class")).
		Part(protocol.NewFileBytes([]byte(class.Content())).WithMetadata("class-name", class.Name())).
		Build()
	resp, err := critic.SendMessage(
		protocol.NewMessageSendParamsBuilder().
			Message(msg).
			Build())
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
	var java string
	var task string
	for _, part := range m.Parts {
		partKind := part.PartKind()
		if partKind == protocol.PartKindText {
			task = part.(*protocol.TextPart).Text
			c.log.Debug("received task: %s", task)
		}
		if partKind == protocol.PartKindFile {
			filePart := part.(*protocol.FilePart)
			content, err := base64.StdEncoding.DecodeString(filePart.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, err
			}
			java = string(content)
			c.log.Debug("received file: %s", content)
		}
	}
	c.log.Info("received messsage #%s, '%s', number of attached files: %d", m.MessageID, task, 1)
	c.log.Info("asking ai to find flaws in the code...")
	replacer := strings.NewReplacer(
		"{{code}}", java,
		"{{imperfections}}", tool.NewCombined(c.tools...).Imperfections(),
	)
	answer, err := c.brain.Ask(replacer.Replace(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from brain: %w", err)
	}
	builder := protocol.NewMessageBuilder().MessageID(m.MessageID)
	suggestions := parseAnswer(answer)
	c.log.Info("found %d possible improvements", len(suggestions))
	for _, suggestion := range suggestions {
		c.log.Debug("suggestion: %s", suggestion)
		builder.Part(protocol.NewText(suggestion))
	}
	res := builder.Build()
	c.log.Debug("sending response: %s", res.MessageID)
	return res, nil
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

func agentCard(port int) *protocol.AgentCard {
	return protocol.Card().
		Name("Critic Agent").
		Description("Critic Description").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("critic-java-code", "Critic Java Code", "Give a reasonable critique on Java code").
		Build()
}
