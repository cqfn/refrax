package domain

import (
	"fmt"
	"strconv"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
	"github.com/google/uuid"
)

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
		case "class":
			c, err := UnmarshalClass(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal class: %w", err)
			}
			classes = append(classes, c)
		case "example":
			e, err := UnmarshalClass(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal example class: %w", err)
			}
			examples = append(examples, e)
		case "suggestion":
			s, err := UnmarshalSuggestion(part)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal suggestion: %w", err)
			}
			suggestions = append(suggestions, s)
		default:
			return nil, fmt.Errorf("unknown part type %s", t)
		}
	}
	job.Classes = classes
	job.Examples = examples
	job.Suggestions = suggestions
	return job, nil
}

func UnmarshalSuggestion(part protocol.Part) (Suggestion, error) {
	if part.PartKind() == protocol.PartKindText {
		suggestion := NewSuggestion(part.(*protocol.TextPart).Text)
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
	return NewClass(
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
		meta: part.Metadata(),
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
			msg.AddPart(MarshalClass(class, "class"))
		}
	}
	if len(j.Suggestions) > 0 {
		for _, suggestion := range j.Suggestions {
			if suggestion == nil {
				continue
			}
			msg.AddPart(MarshalSuggestion(suggestion))
		}
	}
	if len(j.Examples) > 0 {
		for _, example := range j.Examples {
			if example == nil {
				continue
			}
			msg.AddPart(MarshalClass(example, "example"))
		}
	}
	return protocol.NewMessageSendParams().WithMessage(msg)
}

func (d *Description) Marshal() protocol.Part {
	part := protocol.NewText(d.Text)
	for k, v := range d.meta {
		part = part.WithMetadata(k, v)
	}
	return part
}

func MarshalClass(c Class, t string) protocol.Part {
	return protocol.NewFileBytes([]byte(c.Content())).
		WithMetadata("type", t).
		WithMetadata("class-name", c.Name()).
		WithMetadata("class-path", c.Path())
}

func MarshalSuggestion(s Suggestion) protocol.Part {
	return protocol.NewText(s.Text()).WithMetadata("type", "suggestion")
}

// TaskToMsg converts a Task object to a protocol.Message.
func TaskToMsg(task Task) *protocol.Message {
	param, ok := task.Param("max-size")
	if !ok {
		param = "200"
	}
	size, err := strconv.Atoi(param)
	if err != nil {
		panic(err)
	}
	msg := protocol.NewMessage().
		WithMessageID(uuid.NewString()).
		AddPart(protocol.NewText(task.Description()).WithMetadata("max-size", size))
	all := task.Classes()
	for _, class := range all {
		name := class.Name()
		path := class.Path()
		msg = msg.AddPart(protocol.NewFileBytes([]byte(class.Content())).WithMetadata("class-name", name).WithMetadata("class-path", path))
	}
	return msg
}

// RespToClasses converts a protocol.JSONRPCResponse to a slice of Class objects.
func RespToClasses(resp *protocol.JSONRPCResponse) ([]Class, error) {
	parts := resp.Result.(*protocol.Message).Parts
	log.Debug("received %d parts in refactoring response", len(parts))
	classes := make([]Class, 0, len(parts))
	for _, p := range parts {
		kind := p.PartKind()
		if kind == protocol.PartKindFile {
			log.Debug("received file part %v", p)
			classname := p.Metadata()["class-name"]
			path := fmt.Sprintf("%v", p.Metadata()["class-path"])
			bytes := p.(*protocol.FilePart).File.(protocol.FileWithBytes).Bytes
			decoded, err := util.DecodeFile(bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to decode refactored class %s: %w", classname, err)
			}
			class := NewClass(fmt.Sprintf("%v", classname), path, decoded)
			classes = append(classes, class)
		}
	}
	return classes, nil
}

// ClassesToMsg compiles a protocol message from a list of classes.
func ClassesToMsg(classes []Class) *protocol.Message {
	msg := protocol.NewMessage()
	for _, v := range classes {
		msg.AddPart(protocol.NewFileBytes([]byte(v.Content())).WithMetadata("class-name", v.Name()))
	}
	return msg
}

// RespToSuggestions converts a protocol.JSONRPCResponse to a slice of Suggestion objects.
func RespToSuggestions(resp *protocol.JSONRPCResponse) []Suggestion {
	criticMessage := resp.Result.(*protocol.Message)
	suggestions := make([]Suggestion, 0, len(criticMessage.Parts))
	for _, part := range criticMessage.Parts {
		if part.PartKind() == protocol.PartKindText {
			suggestion := NewSuggestion(part.(*protocol.TextPart).Text)
			suggestions = append(suggestions, suggestion)
		}
	}
	return suggestions
}

// RespToClass converts a JSONRPCResponse into a Class, decoding file parts as needed.
func RespToClass(resp *protocol.JSONRPCResponse) (Class, error) {
	part := resp.Result.(*protocol.Message).Parts[0]
	if part.PartKind() == protocol.PartKindFile {
		filePart := part.(*protocol.FilePart)
		decoded, err := util.DecodeFile(filePart.File.(protocol.FileWithBytes).Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to decode file part: %w", err)
		}
		meta := filePart.Metadata()
		name := meta["class-name"]
		if name == nil {
			return nil, fmt.Errorf("file part has no class name metadata")
		}
		path := meta["class-path"]
		if path == nil {
			return nil, fmt.Errorf("file part has no class path metadata")
		}
		log.Debug("decoded class name: %s, path: ", name, path)
		return NewClass(fmt.Sprintf("%v", name), fmt.Sprintf("%v", path), decoded), nil
	}
	return nil, fmt.Errorf("expected file part, got %s", part.PartKind())
}

// MsgToTask parses a task from a protocol message.
func MsgToTask(msg *protocol.Message) (Task, error) {
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
		class := NewClass(fmt.Sprintf("%v", meta["class-name"]), fmt.Sprintf("%v", meta["class-path"]), decoded)
		classes = append(classes, class)
	}
	return &task{
		descr:      descr,
		classes:    classes,
		parameters: params,
	}, nil
}
