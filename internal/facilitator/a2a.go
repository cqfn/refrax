package facilitator

import (
	"fmt"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
)

type a2aTask struct {
	classes []Class
	descr   string
	example Class
	params  map[string]any
}

// Classes implements Task.
func (a *a2aTask) Classes() []Class {
	return a.classes
}

// Description implements Task.
func (a *a2aTask) Description() string {
	return a.descr
}

// Example implements Task.
func (a *a2aTask) Example() Class {
	return a.example
}

// Param implements Task.
func (a *a2aTask) Param(name string) (string, bool) {
	v, ok := a.params[name]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", v), true
}

type a2aSuggestion struct {
	protocol.TextPart
}

// Text implements Suggestion.
func (a *a2aSuggestion) Text() string {
	return a.TextPart.Text
}

type a2aClass struct {
	content string
	name    string
}

// Content implements Class.
func (a *a2aClass) Content() string {
	return a.content
}

// Name implements Class.
func (a *a2aClass) Name() string {
	return a.name
}

// ParseSuggestions parses suggestions from a JSON RPC response.
func ParseSuggestions(resp *protocol.JSONRPCResponse) []Suggestion {
	criticMessage := resp.Result.(protocol.Message)
	suggestions := make([]Suggestion, 0, len(criticMessage.Parts)) // Pre-allocate memory
	for _, part := range criticMessage.Parts {
		suggestion := ParseSuggestion(part)
		suggestions = append(suggestions, suggestion)
	}
	return suggestions
}

// ParseSuggestion parses a single suggestion from a protocol part.
func ParseSuggestion(part protocol.Part) Suggestion {
	if part.PartKind() == protocol.PartKindText {
		return &a2aSuggestion{TextPart: *part.(*protocol.TextPart)}
	}
	return nil
}

// ParseClass parses a class from a JSON RPC response.
func ParseClass(resp *protocol.JSONRPCResponse) (Class, error) {
	part := resp.Result.(protocol.Message).Parts[0]
	if part.PartKind() == protocol.PartKindFile {
		filePart := part.(*protocol.FilePart)
		decoded, err := util.DecodeFile(filePart.File.(protocol.FileWithBytes).Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to decode file part: %w", err)
		}
		meta := filePart.Metadata()
		if meta["class-name"] == nil {
			return nil, fmt.Errorf("file part has no class name metadata")
		}
		log.Debug("decoded class name: %s", meta["class-name"])
		return &a2aClass{
			content: decoded,
			name:    fmt.Sprintf("%v", meta["class-name"]),
		}, nil
	}
	return nil, fmt.Errorf("expected file part, got %s", part.PartKind())
}

// ParseTask parses a task from a protocol message.
func ParseTask(msg *protocol.Message) (Task, error) {
	if len(msg.Parts) == 0 {
		return nil, fmt.Errorf("message has no parts")
	}
	part := msg.Parts[0]
	if part.PartKind() != protocol.PartKindText {
		return nil, fmt.Errorf("expected text part, got %s", part.PartKind())
	}
	descr := part.(*protocol.TextPart).Text
	params := part.Metadata()
	classes := make([]Class, 0)
	var example Class

	for _, v := range msg.Parts[1:] {
		if v.PartKind() != protocol.PartKindFile {
			continue
		}
		filePart := v.(*protocol.FilePart)
		decoded, err := util.DecodeFile(filePart.File.(protocol.FileWithBytes).Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to decode file part: %w", err)
		}
		meta := v.Metadata()
		class := &a2aClass{
			content: decoded,
			name:    fmt.Sprintf("%v", meta["class-name"]),
		}
		if meta["example"] != nil {
			example = class
			continue
		}
		classes = append(classes, class)
	}

	return &a2aTask{
		descr:   descr,
		classes: classes,
		example: example,
		params:  params,
	}, nil
}

// CompileMessage compiles a protocol message from a list of classes.
func CompileMessage(classes []Class) *protocol.Message {
	msg := protocol.NewMessageBuilder()
	for _, v := range classes {
		msg.Part(protocol.NewFileBytes([]byte(v.Content())).WithMetadata("class-name", v.Name()))
	}
	return msg.Build()
}
