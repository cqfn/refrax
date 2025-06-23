package facilitator

import (
	"encoding/base64"
	"fmt"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Facilitator struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
}

func NewFacilitator(ai brain.Brain, port int) (*Facilitator, error) {
	logger := log.NewPrefixed("facilitator: ", log.Default())
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewCustomServer(agentCard(port), port)
	facilitator := &Facilitator{
		server: server,
		brain:  ai,
		log:    logger,
	}
	server.SetHandler(facilitator.think)
	return facilitator, nil
}

func StartServer(ai string, token string, port int) error {
	facilitator, err := NewFacilitator(brain.New(ai, token), port)
	facilitator.log.Debug("created facilitator server: %v", facilitator)
	if err != nil {
		return fmt.Errorf("failed to create facilitator: %w", err)
	}
	facilitator.log.Debug("starting facilitator server on port %d", port)
	return facilitator.Start(make(chan struct{}))
}

func (f *Facilitator) Start(ready chan<- struct{}) error {
	f.log.Info("starting facilitator server...")
	if err := f.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
	f.log.Info("facilitator server has started successfully")
	return nil
}

func (f *Facilitator) Close() error {
	f.log.Info("stopping facilitator server...")
	if err := f.server.Close(); err != nil {
		return fmt.Errorf("failed to stop facilitator server: %w", err)
	}
	f.log.Info("facilitator server stopped successfully")
	return nil
}

func (f *Facilitator) think(m *protocol.Message) (*protocol.Message, error) {
	f.log.Info("received message id:  %s", m.MessageID)
	f.log.Debug("message parts: %v", m.Parts)
	var java string
	for _, part := range m.Parts {
		partKind := part.PartKind()
		if partKind == protocol.PartKindText {
			task := part.(*protocol.TextPart).Text
			f.log.Info("received task: %s", task)
		}
		if partKind == protocol.PartKindFile {
			filePart := part.(*protocol.FilePart)
			f.log.Info("received file: %s", filePart.File)
			content, err := base64.StdEncoding.DecodeString(filePart.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, err
			}
			java = string(content)
			f.log.Debug("file content: %s", content)
		}
	}
	prompt := fmt.Sprintf("Refactor the following Java code to improve its structure and readability:\n```java\n%s\n```", java)
	answer, err := f.brain.Ask(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from brain: %w", err)
	}
	res := protocol.NewMessageBuilder().
		Part(protocol.NewFileBytes([]byte(answer))).
		Build()
	f.log.Debug("sending response: %s", res)
	return &res, nil
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
