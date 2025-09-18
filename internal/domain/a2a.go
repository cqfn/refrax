package domain

import (
	"fmt"

	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
	"github.com/google/uuid"
)

const (
	typeClass      = "class"
	typeSuggestion = "suggestion"
	typeExample    = "example"
)

func UnmarshalArtifacts(msg *protocol.Message) (*Artifacts, error) {
	if len(msg.Parts) == 0 {
		return nil, fmt.Errorf("message has no parts")
	}
	artifacts := &Artifacts{}
	descr, err := UnmarshalDescription(msg.Parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal description: %w", err)
	}
	artifacts.Descr = descr
	classes := make([]Class, 0)
	suggestions := make([]Suggestion, 0)
	for _, part := range msg.Parts[1:] {
		metas := part.Metadata()
		t := metas["type"]
		switch t {
		case typeClass:
			c, err := UnmarshalClass(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal class: %w", err)
			}
			classes = append(classes, c)
		case typeSuggestion:
			s, err := UnmarshalSuggestion(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal suggestion: %w", err)
			}
			suggestions = append(suggestions, *s)
		default:
			return nil, fmt.Errorf("unknown part type %s", t)
		}
	}
	artifacts.Classes = classes
	artifacts.Suggestions = suggestions
	return artifacts, nil
}

func UnmarshalJob(msg *protocol.Message) (*Job, error) {
	if len(msg.Parts) == 0 {
		return nil, fmt.Errorf("message has no parts")
	}
	job := &Job{}
	descr, err := UnmarshalDescription(msg.Parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal description: %w", err)
	}
	job.Descr = descr
	classes := make([]Class, 0)
	examples := make([]Class, 0)
	suggestions := make([]Suggestion, 0)
	for _, part := range msg.Parts[1:] {
		metas := part.Metadata()
		t := metas["type"]
		switch t {
		case typeClass:
			c, err := UnmarshalClass(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal class: %w", err)
			}
			classes = append(classes, c)
		case typeExample:
			e, err := UnmarshalClass(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal example class: %w", err)
			}
			examples = append(examples, e)
		case typeSuggestion:
			s, err := UnmarshalSuggestion(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal suggestion: %w", err)
			}
			suggestions = append(suggestions, *s)
		default:
			return nil, fmt.Errorf("unknown part type %s", t)
		}
	}
	job.Classes = classes
	job.Examples = examples
	job.Suggestions = suggestions
	return job, nil
}

func UnmarshalOldSuggestion(part protocol.Part) (OldSuggestion, error) {
	if part.PartKind() == protocol.PartKindText {
		suggestion := NewOldSuggestion(part.(*protocol.TextPart).Text)
		return suggestion, nil
	}
	return nil, fmt.Errorf("expected text part for suggestion, got %s", part.PartKind())
}

func UnmarshalSuggestion(part protocol.Part) (*Suggestion, error) {
	if part.PartKind() == protocol.PartKindText {
		tp := part.(*protocol.TextPart)
		text := tp.Text
		path, ok := tp.Metadata()["class-path"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid class-path metadata in suggestion part")
		}
		suggestion := &Suggestion{
			ClassPath: path,
			Text:      text,
		}
		return suggestion, nil
	}
	return nil, fmt.Errorf("expected text part for suggestion, got %s", part.PartKind())
}

func UnmarshalClass(part protocol.Part) (Class, error) {
	if part.PartKind() != protocol.PartKindFile {
		return nil, fmt.Errorf("expected file part, got %s", part.PartKind())
	}
	f := part.(*protocol.FilePart)
	decoded, err := util.DecodeFile(f.File.(protocol.FileWithBytes).Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file part: %w", err)
	}
	meta := f.Metadata()
	return NewInMemoryClass(
		fmt.Sprintf("%v", meta["class-name"]),
		fmt.Sprintf("%v", meta["class-path"]),
		decoded,
	), nil
}

func UnmarshalDescription(part protocol.Part) (*Description, error) {
	if part.PartKind() != protocol.PartKindText {
		return nil, fmt.Errorf("expected text part for description, got %s", part.PartKind())
	}
	descr := part.(*protocol.TextPart).Text
	res := Description{
		Text: descr,
		Meta: part.Metadata(),
	}
	return &res, nil
}

func (j *Job) Marshal() *protocol.MessageSendParams {
	msg := protocol.NewMessage().WithMessageID(uuid.NewString())
	if j.Descr != nil {
		msg.AddPart(j.Descr.Marshal())
	}
	if len(j.Classes) > 0 {
		for _, class := range j.Classes {
			if class == nil {
				continue
			}
			msg.AddPart(MarshalClass(class, typeClass))
		}
	}
	if len(j.Suggestions) > 0 {
		for _, suggestion := range j.Suggestions {
			msg.AddPart(suggestion.Marshal())
		}
	}
	if len(j.Examples) > 0 {
		for _, example := range j.Examples {
			if example == nil {
				continue
			}
			msg.AddPart(MarshalClass(example, typeExample))
		}
	}
	return protocol.NewMessageSendParams().WithMessage(msg)
}

func (a *Artifacts) Marshal() *protocol.MessageSendParams {
	msg := protocol.NewMessage().WithMessageID(uuid.NewString())
	if a.Descr != nil {
		msg.AddPart(a.Descr.Marshal())
	}
	if len(a.Classes) > 0 {
		for _, class := range a.Classes {
			if class == nil {
				continue
			}
			msg.AddPart(MarshalClass(class, typeClass))
		}
	}
	if len(a.Suggestions) > 0 {
		for _, suggestion := range a.Suggestions {
			msg.AddPart(suggestion.Marshal())
		}
	}
	return protocol.NewMessageSendParams().WithMessage(msg)
}

func (d *Description) Marshal() protocol.Part {
	part := protocol.NewText(d.Text)
	for k, v := range d.Meta {
		part = part.WithMetadata(k, v)
	}
	return part
}

func (s *Suggestion) Marshal() protocol.Part {
	part := protocol.NewText(s.Text).
		WithMetadata("class-path", s.ClassPath).
		WithMetadata("type", typeSuggestion)
	return part
}

func MarshalClass(c Class, t string) protocol.Part {
	return protocol.NewFileBytes([]byte(c.Content())).
		WithMetadata("type", t).
		WithMetadata("class-name", c.Name()).
		WithMetadata("class-path", c.Path())
}

func MarshalSuggestion(s OldSuggestion) protocol.Part {
	return protocol.NewText(s.Text()).WithMetadata("type", typeSuggestion)
}
