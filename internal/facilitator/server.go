package facilitator

import (
	"encoding/base64"
	"fmt"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Facilitator struct {
	server protocol.Server
}

func NewFacilitator(ai string, port int) (*Facilitator, error) {
	log.Debug(prefixed("preparing server on port %d with ai provider %s"), port, ai)
	server, error := protocol.NewCustomServer(agentCard(port), protocol.LogRequest, port)
	if error != nil {
		return nil, fmt.Errorf("failed to create A2A server: %w", error)
	}
	res := &Facilitator{
		server: server,
	}
	server.SetHandler(res.brainLess)
	return res, nil
}

func StartServer(ai string, port int) error {
	log.Debug(prefixed("starting facilitator server on port %d with ai provider %s"), port, ai)
	facilitator, err := NewFacilitator(ai, port)
	if err != nil {
		return fmt.Errorf("failed to create facilitator: %w", err)
	}
	return facilitator.Start(make(chan struct{}))
}

func (f *Facilitator) Start(ready chan<- struct{}) error {
	log.Info(prefixed("starting facilitator server"))
	if err := f.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
	log.Info(prefixed("facilitator server started successfully"))
	return nil
}

func (f *Facilitator) Close() error {
	log.Info(prefixed("stopping facilitator server..."))
	if err := f.server.Close(); err != nil {
		return fmt.Errorf("failed to stop facilitator server: %w", err)
	}
	log.Info(prefixed("facilitator server stopped successfully"))
	return nil
}

func (f *Facilitator) brainLess(m *protocol.Message) (*protocol.Message, error) {
	log.Info(prefixed("received message id:  %s"), m.MessageID)
	log.Debug(prefixed("message parts: %v"), m.Parts)
	for _, part := range m.Parts {
		partKind := part.PartKind()
		if partKind == protocol.PartKindText {
			task := part.(*protocol.TextPart).Text
			log.Info(prefixed("received task: %s"), task)
		}
		if partKind == protocol.PartKindFile {
			filePart := part.(*protocol.FilePart)
			log.Info(prefixed("received file: %s"), filePart.File)
			content, err := base64.StdEncoding.DecodeString(filePart.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, err
			}
			log.Debug(prefixed("file content: %s"), content)
		}
	}
	res := protocol.NewMessageBuilder().
		Part(protocol.NewFileBytes([]byte(refactored()))).
		Build()
	log.Debug(prefixed("sending response: %s"), res)
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

func refactored() string {
	return "public class Main {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello, World\");\n\t}\n"
}

func prefixed(template string) string {
	return fmt.Sprintf("facilitator: %s", template)
}
