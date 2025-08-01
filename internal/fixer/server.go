// Package fixer provides functionality for fixing Java code based on suggestions.
package fixer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
	"github.com/google/uuid"
)

// Fixer is a server that fixes Java code based on suggestions provided.
type Fixer struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
	port   int
}

const prompt = "Fix '%s' code based on the listed suggestions.\n\n" +
	"Code:\n" +
	"```java\n%s\n\n```" +
	"\n\nSuggestions:\n" +
	"```suggestion\n%s\n\n```" +
	"Do not rename class names.\n" +
	"Return only the corrected Java code.\n" +
	"Do not include explanations, comments, or any extra text.\n\n"

const exampleNote = "This is an example of another refactored class. The same refactoring has been applied here." +
	"```java\n%s\n```"

// NewFixer creates a new Fixer instance with the provided AI brain and port.
func NewFixer(ai brain.Brain, port int) *Fixer {
	logger := log.NewPrefixed("fixer", log.NewColored(log.Default(), log.Magenta))
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewServer(agentCard(port), port)
	fixer := &Fixer{
		server: server,
		brain:  ai,
		log:    logger,
		port:   port,
	}
	server.MsgHandler(fixer.think)
	fixer.log.Debug("preparing the Fixer server on port %d with ai provider %s", port, ai)
	return fixer
}

// Fix applies the given suggestions to the provided class and returns the modified class or an error.
// It communicates with an external fixer service to perform the modifications.
func (f *Fixer) Fix(class domain.Class, suggestions []domain.Suggestion, example domain.Class) (domain.Class, error) {
	address := fmt.Sprintf("http://localhost:%d", f.port)
	f.log.Info("asking fixer (%s) to apply suggestions...", address)
	fixer := protocol.NewClient(address)
	msg := protocol.NewMessage().
		WithMessageID(uuid.NewString()).
		AddPart(protocol.NewText("apply all the following suggestions"))
	for _, suggestion := range suggestions {
		msg.AddPart(protocol.NewText(suggestion.Text()).WithMetadata("suggestion", true))
	}
	if example != nil {
		msg.AddPart(protocol.NewFileBytes([]byte(example.Content())).
			WithMetadata("class-name", example.Name()).
			WithMetadata("example", true))
	}
	file := protocol.NewFileBytes([]byte(class.Content()))
	msg = msg.AddPart(file.WithMetadata("class-name", class.Name()))
	resp, err := fixer.SendMessage(protocol.NewMessageSendParams().WithMessage(msg))
	if err != nil {
		return nil, fmt.Errorf("failed to send message to fixer: %w", err)
	}
	return domain.RespToClass(resp)
}

// ListenAndServe begins the Fixer server and signals readiness through the provided channel.
func (f *Fixer) ListenAndServe() error {
	f.log.Info("starting fixer server on port %d...", f.port)
	var err error
	if err = f.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start fixer server: %w", err)
	}
	return err
}

// Shutdown gracefully stops the Fixer server.
func (f *Fixer) Shutdown() error {
	f.log.Info("stopping fixer server...")
	if err := f.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop fixer server: %w", err)
	}
	f.log.Info("fixer server stopped successfully")
	return nil
}

// Ready returns a channel that signals when the Fixer server is ready to accept requests.
func (f *Fixer) Ready() <-chan bool {
	return f.server.Ready()
}

// Handler sets the handler function for processing requests on the Fixer server.
func (f *Fixer) Handler(hander protocol.Handler) {
	f.server.Handler(hander)
}

func (f *Fixer) think(ctx context.Context, m *protocol.Message) (*protocol.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled: %w", ctx.Err())
	case res := <-f.thinkChan(m):
		return res.msg, res.err
	}
}

type thought struct {
	msg *protocol.Message
	err error
}

func (f *Fixer) thinkChan(m *protocol.Message) <-chan thought {
	res := make(chan thought, 1)
	go func() {
		msg, err := f.thinkLong(m)
		res <- thought{
			msg, err,
		}
		close(res)
	}()
	return res
}

func (f *Fixer) thinkLong(m *protocol.Message) (*protocol.Message, error) {
	f.log.Info("received message: #%s", m.MessageID)
	f.log.Info("trying to fix Java code...")
	var code string
	var suggestions []string
	var class string
	var example string
	for _, part := range m.Parts {
		if part.PartKind() == protocol.PartKindText {
			msg := part.(*protocol.TextPart).Text
			if part.Metadata()["suggestion"] != nil {
				suggestions = append(suggestions, msg)
			}
		} else if part.PartKind() == protocol.PartKindFile {
			file := part.(*protocol.FilePart)
			var err error
			content, err := util.DecodeFile(file.File.(protocol.FileWithBytes).Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to decode code part: %w", err)
			}
			if part.Metadata()["example"] != nil {
				example = content
			} else {
				code = content
				name := file.Metadata()["class-name"]
				class = fmt.Sprintf("%v", name)
			}

		}
	}
	all := strings.Join(suggestions, "\n")
	question := fmt.Sprintf(prompt, class, code, all)
	if example != "" {
		question += "\n\n" + fmt.Sprintf(exampleNote, example)
	}
	f.log.Info("asking AI to fix java code...")
	f.log.Debug("asking AI: %s", question)
	answer, err := f.brain.Ask(question)
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from AI: %w", err)
	}
	f.log.Debug("received answer from AI: %s", answer)
	f.log.Info("AI provided a fix for the Java code, sending response back...")
	return protocol.NewMessage().
		WithMessageID(m.MessageID).
		AddPart(protocol.NewFileBytes([]byte(clean(answer))).WithMetadata("class-name", class)), nil
}

func clean(answer string) string {
	answer = strings.ReplaceAll(answer, "```java", "")
	return strings.ReplaceAll(answer, "```", "")
}

func agentCard(port int) *protocol.AgentCard {
	return protocol.NewAgentCard().
		WithName("Fixer Agent").
		WithDescription("Fixer Description").
		WithURL(fmt.Sprintf("http://localhost:%d", port)).
		WithVersion("0.0.1").
		AddSkill("fix-java-code", "Fix Java Code", "Fix a Java code based on suggestions")
}
