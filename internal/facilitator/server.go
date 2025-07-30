// Package facilitator provides functionality for facilitating interactions between
// critic and fixer agents in a code refactoring process.
package facilitator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
)

// Facilitator facilitates communication between the critic and fixer agents.
type Facilitator struct {
	server     protocol.Server
	brain      brain.Brain
	log        log.Logger
	port       int
	criticPort int
	fixerPort  int
}

// NewFacilitator creates a new instance of Facilitator to manage communication between agents.
func NewFacilitator(ai brain.Brain, port, criticPort, fixerPort int) *Facilitator {
	logger := log.NewPrefixed("facilitator", log.Default())
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewServer(agentCard(port), port)
	facilitator := &Facilitator{
		server:     server,
		brain:      ai,
		log:        logger,
		criticPort: criticPort,
		fixerPort:  fixerPort,
		port:       port,
	}
	server.MsgHandler(facilitator.think)
	return facilitator
}

// ListenAndServe starts the facilitator server and prepares it for handling requests.
func (f *Facilitator) ListenAndServe() error {
	f.log.Info("starting facilitator server on port %d...", f.port)
	var err error
	if err = f.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start facilitator server: %w", err)
	}
	return err
}

// Shutdown stops the facilitator server and releases resources.
func (f *Facilitator) Shutdown() error {
	f.log.Info("stopping facilitator server...")
	if err := f.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop facilitator server: %w", err)
	}
	f.log.Info("facilitator server stopped successfully")
	return nil
}

// Ready returns a channel that signals when the facilitator server is ready to accept requests.
func (f *Facilitator) Ready() <-chan bool {
	return f.server.Ready()
}

// Handler sets the message handler for the facilitator server.
func (f *Facilitator) Handler(handler protocol.Handler) {
	f.server.Handler(handler)
}

func (f *Facilitator) think(ctx context.Context, m *protocol.Message) (*protocol.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return f.thinkLong(m)
	}
}

func (f *Facilitator) thinkLong(m *protocol.Message) (*protocol.Message, error) {
	if f.isRefactoringTask(m) {
		return f.refactor(m)
	}
	f.log.Warn("received a message that is not related to refactoring. ignoring.")
	return nil, fmt.Errorf("received a message that is not related to refactoring")
}

func (f *Facilitator) isRefactoringTask(m *protocol.Message) bool {
	res := false
	if len(m.Parts) > 0 {
		for _, part := range m.Parts {
			if part.PartKind() == protocol.PartKindText {
				task := part.(*protocol.TextPart).Text
				if strings.Contains(task, "refactor the project") {
					res = true
					break
				}
			}
		}
	}
	if !res {
		f.log.Warn("received a task that is not related to refactoring. ignoring.")
	}
	return res
}

func (f *Facilitator) refactor(m *protocol.Message) (*protocol.Message, error) {
	maxSize := maxSizeParam(m)
	f.log.Info("received messsage %s, number of attached files: %d, max-size: %d", m.MessageID, nfiles(m), maxSize)
	response := protocol.NewMessageBuilder()
	changed := 0
	var example string
	for _, part := range m.Parts {
		var file *protocol.FilePart
		var class string
		partKind := part.PartKind()
		if partKind == protocol.PartKindFile {
			class = part.Metadata()["class-name"].(string)
			filePart := part.(*protocol.FilePart)
			file = filePart
			content, err := util.DecodeFile(filePart.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed decode file: %w", err)
			}

			f.log.Debug("received file: %s", content)
			if changed >= maxSize {
				f.log.Info("refactoring class %s would exceed max size %d, skipping refactoring", class, maxSize)
				response.Part(file.WithMetadata("class-name", class).WithMetadata("refactor-status", "skipped"))
				continue
			}
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

			fixed, err := f.AskFixer(m.MessageID, suggestions, class, file, example)
			if err != nil {
				return nil, fmt.Errorf("failed to ask fixer: %w", err)
			}
			filePartResult := fixed.Result.(protocol.Message).Parts[0].(*protocol.FilePart)
			after, err := util.DecodeFile(filePartResult.File.(protocol.FileWithBytes).Bytes)
			example = after
			if err != nil {
				return nil, fmt.Errorf("failed to decode fixed file: %w", err)
			}
			changed += util.Diff(content, after)
			f.log.Info("total number of fixed lines: %s", changed)
			f.log.Info("received fixed file from fixer, sending final response...")
			response = response.Part(filePartResult.WithMetadata("class-name", class))
		}
	}
	res := response.Build()
	f.log.Debug("sending response: %s", res.MessageID)
	return res, nil
}

func maxSizeParam(m *protocol.Message) int {
	res := 0
	for _, v := range m.Parts {
		size := v.Metadata()["max-size"]
		if size != nil {
			return int(size.(float64))
		}
	}
	return res
}

func nfiles(m *protocol.Message) int {
	res := 0
	for _, v := range m.Parts {
		if v.PartKind() == protocol.PartKindFile {
			res++
		}
	}
	return res
}

// AskCritic sends a file to the critic agent for linting and analysis.
func (f *Facilitator) AskCritic(id string, file *protocol.FilePart) (*protocol.JSONRPCResponse, error) {
	address := fmt.Sprintf("http://localhost:%d", f.criticPort)
	f.log.Info("asking critic (%s) to lint the class...", address)
	f.log.Debug("message id: %s, file: %s", id, file.File.(protocol.FileWithBytes))
	critic := protocol.NewClient(address)
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

// AskFixer sends the suggestions and file to the fixer agent for processing.
func (f *Facilitator) AskFixer(id string, suggestions []string, class string, file *protocol.FilePart, example string) (*protocol.JSONRPCResponse, error) {
	address := fmt.Sprintf("http://localhost:%d", f.fixerPort)
	log.Debug("asking fixer (%s) to apply suggestions...", address)
	fixer := protocol.NewClient(address)
	builder := protocol.NewMessageBuilder().
		MessageID(id).
		Part(protocol.NewText("apply all the following suggestions"))
	for _, suggestion := range suggestions {
		builder.Part(protocol.NewText(suggestion).WithMetadata("suggestion", true))
	}
	if example != "" {
		builder.Part(protocol.NewText(example).WithMetadata("example", true))
	}
	msg := builder.Part(file.WithMetadata("class-name", class)).Build()
	return fixer.SendMessage(protocol.NewMessageSendParamsBuilder().Message(msg).Build())
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

func className(task string) string {
	begin := strings.Index(task, "'") + 1
	end := begin + strings.Index(task[begin:], "'")
	if begin >= end || begin < 0 || end < 0 {
		return ""
	}
	return task[begin:end]
}
