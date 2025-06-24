package fixer

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Fixer struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
	port   int
}

const prompt = `Fix the following Java code based on the listed suggestions.

Code:
%s

Suggestions:
%s

Do not rename class names.
Return only the corrected Java code.
Do not include explanations, comments, or any extra text.`

func NewFixer(ai brain.Brain, port int) *Fixer {
	logger := log.NewPrefixed("fixer", log.Default())
	server := protocol.NewCustomServer(agentCard(port), port)
	fixer := &Fixer{
		server: server,
		brain:  ai,
		log:    logger,
		port:   port,
	}
	server.SetHandler(fixer.think)
	fixer.log.Debug("preparing the Fixer server on port %d with ai provider %s", port, ai)
	return fixer
}

func (c *Fixer) Start(ready chan<- struct{}) error {
	c.log.Info("starting fixer server on port %d...", c.port)
	if err := c.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start fixer server: %w", err)
	}
	return nil
}

func (c *Fixer) Close() error {
	c.log.Info("stopping fixer server...")
	if err := c.server.Close(); err != nil {
		return fmt.Errorf("failed to stop fixer server: %w", err)
	}
	c.log.Info("fixer server stopped successfully")
	return nil
}

func (c *Fixer) think(m *protocol.Message) (*protocol.Message, error) {
	c.log.Info("received message: #%s", m.MessageID)
	c.log.Info("trying to fix Java code...")
	var code string
	var suggestions []string
	for _, part := range m.Parts {
		if part.PartKind() == protocol.PartKindText {
			suggestion := part.(*protocol.TextPart).Text
			if suggestion != "" {
				suggestions = append(suggestions, suggestion)
			}
		} else if part.PartKind() == protocol.PartKindFile {
			codePart := part.(*protocol.FilePart)
			codeBytes, err := base64.StdEncoding.DecodeString(codePart.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to decode code part: %w", err)
			}
			code = string(codeBytes)
		}
	}
	all := strings.Join(suggestions, "\n")
	question := fmt.Sprintf(prompt, code, all)
	c.log.Info("asking AI to fix java code...")
	c.log.Debug("asking AI: %s", question)
	answer, err := c.brain.Ask(question)
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from AI: %w", err)
	}
	c.log.Debug("received answer from AI: %s", answer)
	c.log.Info("AI provided a fix for the Java code, sending response back...")
	message := protocol.NewMessageBuilder().
		MessageID(m.MessageID).
		Part(protocol.NewFileBytes([]byte(clean(answer)))).
		Build()
	return &message, nil
}

func clean(answer string) string {
	answer = strings.ReplaceAll(answer, "```java", "")
	return strings.ReplaceAll(answer, "```", "")
}

func agentCard(port int) protocol.AgentCard {
	return protocol.Card().
		Name("Fixer Agent").
		Description("Fixer Description").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("fix-java-code", "Fix Java Code", "Fix a Java code based on suggestions").
		Build()
}
