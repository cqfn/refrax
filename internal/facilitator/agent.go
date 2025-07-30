package facilitator

import (
	"fmt"
	"strconv"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
	"github.com/google/uuid"
)

type agent struct {
	brain      brain.Brain
	log        log.Logger
	criticPort int
	fixerPort  int
}

// Task represents a unit of work that contains classes and associated parameters.
type Task interface {
	Description() string
	Classes() []Class
	Example() Class
	Param(name string) (string, bool)
}

// Class represents a code or text entity with a name and content.
type Class interface {
	Name() string
	Content() string
}

// Suggestion represents a proposed improvement or fix for a class.
type Suggestion interface {
	Text() string
}

func limit(t Task) (int, error) {
	size, ok := t.Param("max-size")
	if !ok {
		size = "200"
	}
	return strconv.Atoi(size)
}

func (a *agent) AskCritic(class Class) ([]Suggestion, error) {
	address := fmt.Sprintf("http://localhost:%d", a.criticPort)
	a.log.Info("asking critic (%s) to lint the class...", address)
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
	return ParseSuggestions(resp), nil
}

func (a *agent) AskFixer(class Class, suggestions []Suggestion, example Class) (Class, error) {
	address := fmt.Sprintf("http://localhost:%d", a.fixerPort)
	a.log.Info("asking fixer (%s) to apply suggestions...", address)
	fixer := protocol.NewClient(address)
	builder := protocol.NewMessageBuilder().
		MessageID(uuid.NewString()).
		Part(protocol.NewText("apply all the following suggestions"))
	for _, suggestion := range suggestions {
		builder.Part(protocol.NewText(suggestion.Text()).WithMetadata("suggestion", true))
	}
	if example != nil {
		builder.Part(protocol.NewFileBytes([]byte(example.Content())).
			WithMetadata("class-name", example.Name()).
			WithMetadata("example", true))
	}
	file := protocol.NewFileBytes([]byte(class.Content()))
	msg := builder.Part(file.WithMetadata("class-name", class.Name())).Build()
	resp, err := fixer.SendMessage(protocol.NewMessageSendParamsBuilder().Message(msg).Build())
	if err != nil {
		return nil, fmt.Errorf("failed to send message to fixer: %w", err)
	}
	return ParseClass(resp)
}

func (a *agent) refactor(t Task) ([]Class, error) {
	size, err := limit(t)
	if err != nil {
		return nil, fmt.Errorf("failed to get max size limit: %w", err)
	}
	a.log.Info("received request for refactoring, number of attached files: %d, max-size: %d", len(t.Classes()), size)
	refactored := make([]Class, 0, len(t.Classes()))
	var example Class
	changed := 0
	for _, class := range t.Classes() {
		a.log.Info("received class for refactoring: %q", class.Name())
		if changed >= size {
			a.log.Warn("refactoring class %s would exceed max size %d, skipping refactoring", class.Name(), size)
			refactored = append(refactored, class)
			continue
		}
		suggestions, err := a.AskCritic(class)
		if err != nil {
			return nil, fmt.Errorf("failed to ask critic: %w", err)
		}
		a.log.Info("received %d suggestions from critic", len(suggestions))
		modified, err := a.AskFixer(class, suggestions, example)
		if err != nil {
			return nil, fmt.Errorf("failed to ask fixer: %w", err)
		}
		a.log.Info("fixed class %s, changed content", modified.Name())
		refactored = append(refactored, modified)
		diff := util.Diff(class.Content(), modified.Content())
		changed += diff
	}
	return refactored, nil
}
