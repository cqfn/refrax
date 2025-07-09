package facilitator

import (
	"encoding/base64"
	"fmt"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type Facilitator struct {
	server     protocol.Server
	brain      brain.Brain
	log        log.Logger
	port       int
	criticPort int
	fixerPort  int
}

func NewFacilitator(ai brain.Brain, port, criticPort, fixerPort int) *Facilitator {
	logger := log.NewPrefixed("facilitator", log.Default())
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewCustomServer(agentCard(port), port)
	facilitator := &Facilitator{
		server:     server,
		brain:      ai,
		log:        logger,
		criticPort: criticPort,
		fixerPort:  fixerPort,
		port:       port,
	}
	server.SetHandler(facilitator.think)
	return facilitator
}

func (f *Facilitator) Start(ready chan<- struct{}) error {
	f.log.Info("starting facilitator server on port %d...", f.port)
	if err := f.server.Start(ready); err != nil {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
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
	f.log.Debug("received message: #%s", m.MessageID)
	var file *protocol.FilePart
	var task string
	for _, part := range m.Parts {
		partKind := part.PartKind()
		if partKind == protocol.PartKindText {
			task = part.(*protocol.TextPart).Text
			f.log.Debug("received task: %s", task)
		}
		if partKind == protocol.PartKindFile {
			filePart := part.(*protocol.FilePart)
			file = filePart
			content, err := base64.StdEncoding.DecodeString(filePart.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, err
			}
			f.log.Debug("received file: %s", content)
		}
	}
	f.log.Info("received messsage #%s, '%s', number of attached files: %d", m.MessageID, task, 1)
	criticResp, err := f.AskCritic(m.MessageID, file)
	if err != nil {
		return nil, fmt.Errorf("failed to ask critic: %w", err)
	}

	criticMessage := criticResp.Result.(protocol.Message)
	var suggestions []string
	for _, part := range criticMessage.Parts {
		if part.PartKind() == protocol.PartKindText {
			suggestions = append(suggestions, part.(*protocol.TextPart).Text)
			f.log.Debug("received suggestion: %s", part.(*protocol.TextPart).Text)
		}
	}
	f.log.Info("received %d suggestions from critic", len(suggestions))

	fixed, err := f.AskFixer(m.MessageID, suggestions, file)
	if err != nil {
		return nil, fmt.Errorf("failed to ask fixer: %w", err)
	}
	filePartResult := fixed.Result.(protocol.Message).Parts[0].(*protocol.FilePart)
	f.log.Info("received fixed file from fixer, sending final response...")
	res := protocol.NewMessageBuilder().
		Part(filePartResult).
		Build()
	f.log.Debug("sending response: %s", res)
	return res, nil
}

func (f *Facilitator) AskFixer(id string, suggestions []string, file *protocol.FilePart) (*protocol.JSONRPCResponse, error) {
	address := fmt.Sprintf("http://localhost:%d", f.fixerPort)
	log.Debug("asking fixer (%s) to apply suggestions...", address)
	fixer := protocol.NewCustomClient(address)
	builder := protocol.NewMessageBuilder().
		MessageID(id).
		Part(protocol.NewText("apply all the following suggestions"))
	for _, suggestion := range suggestions {
		builder.Part(protocol.NewText(suggestion))
	}
	msg := builder.Part(file).Build()
	return fixer.SendMessage(protocol.NewMessageSendParamsBuilder().Message(msg).Build())
}

func (f *Facilitator) AskCritic(id string, file *protocol.FilePart) (*protocol.JSONRPCResponse, error) {
	address := fmt.Sprintf("http://localhost:%d", f.criticPort)
	f.log.Info("asking critic (%s) to lint the class...", address)
	f.log.Debug("message id: %s, file: %s", id, file.File.(protocol.FileWithBytes))
	critic := protocol.NewCustomClient(address)
	msg := protocol.NewMessageBuilder().
		MessageID(id).
		Part(protocol.NewText("lint class")).
		Part(file).
		Build()
	return critic.SendMessage(
		protocol.NewMessageSendParamsBuilder().
			Message(msg).
			Build(),
	)
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.Card().
		Name("Facilitator Agent").
		Description("An agent that facilitates talk between critic and fixer").
		URL(fmt.Sprintf("http://localhost:%d", port)).
		Version("0.0.1").
		Skill("facilitate-discussion", "Refactor Java Projects", "Facilitate discussion on code refactoring").
		Build()
}
