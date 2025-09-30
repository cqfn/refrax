package fixer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/prompts"
	"github.com/cqfn/refrax/internal/protocol"
)

// Fixer is a server that fixes Java code based on suggestions provided.
type Fixer struct {
	server protocol.Server
	brain  brain.Brain
	log    log.Logger
	port   int
}

// promptData holds the data to be injected into the prompt template.
type promptData struct {
	FilePath    string
	Code        string
	Suggestions []domain.Suggestion
}

// NewFixer creates a new Fixer instance with the provided AI brain and port.
func NewFixer(ai brain.Brain, port int, colorless bool) *Fixer {
	logger := log.New("fixer", log.Magenta, colorless)
	logger.Debug("preparing server on port %d with ai provider %s", port, ai)
	server := protocol.NewServer(agentCard(port), port)
	fixer := &Fixer{
		server: server,
		brain:  ai,
		log:    logger,
		port:   port,
	}
	server.MsgHandler(fixer.think)
	fixer.log.Debug("Preparing the Fixer server on port %d with ai provider %s", port, ai)
	return fixer
}

// Fix applies the given suggestions to the provided class and returns the modified class or an error.
// It communicates with an external fixer service to perform the modifications.
func (f *Fixer) Fix(job *domain.Job) (*domain.Artifacts, error) {
	address := fmt.Sprintf("http://localhost:%d", f.port)
	f.log.Debug("Asking fixer (%s) to apply suggestions...", address)
	fixer := protocol.NewClient(address)
	resp, err := fixer.SendMessage(job.Marshal())
	if err != nil {
		return nil, fmt.Errorf("failed to send message to fixer: %w", err)
	}
	return domain.UnmarshalArtifacts(resp.Result.(*protocol.Message))
}

// ListenAndServe begins the Fixer server and signals readiness through the provided channel.
func (f *Fixer) ListenAndServe() error {
	f.log.Info("Starting fixer server on port %d...", f.port)
	var err error
	if err = f.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start fixer server: %w", err)
	}
	return err
}

// Shutdown gracefully stops the Fixer server.
func (f *Fixer) Shutdown() error {
	f.log.Info("Stopping fixer server...")
	if err := f.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop fixer server: %w", err)
	}
	f.log.Info("Fixer server stopped successfully")
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
	f.log.Info("Received message: #%s", m.MessageID)
	job, err := domain.UnmarshalJob(m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}
	var code string
	var class string
	var path string
	code = job.Classes[0].Content()
	class = job.Classes[0].Name()
	path = job.Classes[0].Path()
	f.log.Info("Trying to fix the %q class...", class)
	prompt := prompts.User{
		Data: promptData{
			FilePath:    path,
			Code:        code,
			Suggestions: job.Suggestions,
		},
		Name: "fixer/fix.md.tmpl",
	}
	question := prompt.String()
	f.log.Debug("Asking the brain to fix the Java code...")
	answer, err := f.brain.Ask(question)
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from AI: %w", err)
	}
	f.log.Debug("Received answer from AI: %s", answer)
	f.log.Info("AI provided a fix for the Java code, sending response back...")
	res := &domain.Artifacts{
		Descr: &domain.Description{
			Text: fmt.Sprintf("Fix for class %s", class),
		},
		Classes: []domain.Class{
			domain.NewInMemoryClass(class, path, clean(answer)),
		},
	}
	return res.Marshal().Message, nil
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
