package critic

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/aibolit"
)

type Critic struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
	port   int
	imperfections[] byte
}

const prompt = `Analyze the following Java code:

{{code}}

Identify possible improvements or flaws such as:

* variables that can be inlined or removed without changing functionality,
* unnecessary comments inside methods,
* redundant code,
* non-idiomatic patterns.

Keep in mind the following imperfections with Java code, identified by automated static analysis system:

{{imperfections}}

Respond with a plain list of suggestions, one per line. Do not include any explanations, summaries, or extra text.
`

func NewCritic(ai brain.Brain, port int, imperfections[] byte) *Critic {
	logger := log.NewPrefixed("critic", log.Default())
	server := protocol.NewCustomServer(agentCard(port), port)
	critic := &Critic{
		server: server,
		brain:  ai,
		log:    logger,
		port:   port,
		imperfections: imperfections,
	}
	server.SetHandler(critic.think)
	critic.log.Debug("preparing the Critic server on port %d with ai provider %s", port, ai)
	return critic
}

func (c *Critic) Start(ready chan<- struct{}) error {
	c.log.Info("starting critic server on port %d...", c.port)
	if err := c.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start critic server: %w", err)
	}
	return nil
}

func (c *Critic) Close() error {
	c.log.Info("stopping critic server...")
	if err := c.server.Close(); err != nil {
		return fmt.Errorf("failed to stop critic server: %w", err)
	}
	c.log.Info("critic server stopped successfully")
	return nil
}

func (c *Critic) think(m *protocol.Message) (*protocol.Message, error) {
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
	 	"{{imperfections}}", aibolit.NewAibolitResponse(string(c.imperfections)).Sanitized(),
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
	c.log.Debug("sending response: %s", res)
	return &res, nil
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

func agentCard(port int) protocol.AgentCard {
	return protocol.Card().
		Name("Critic Agent").
		Description("Critic Description").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("critic-java-code", "Critic Java Code", "Give a reasonable critique on Java code").
		Build()
}
