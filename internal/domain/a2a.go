// Package domain provides functionality for converting tasks, responses, and classes
// to and from protocol messages used in the refactoring process.
package domain

import (
	"fmt"
	"strconv"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/util"
	"github.com/google/uuid"
)

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
		msg = msg.AddPart(protocol.NewFileBytes([]byte(class.Content())).WithMetadata("class-name", name))
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
			bytes := p.(*protocol.FilePart).File.(protocol.FileWithBytes).Bytes
			decoded, err := util.DecodeFile(bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to decode refactored class %s: %w", classname, err)
			}
			class := NewClass(fmt.Sprintf("%v", classname), decoded)
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
		if meta["class-name"] == nil {
			return nil, fmt.Errorf("file part has no class name metadata")
		}
		log.Debug("decoded class name: %s", meta["class-name"])
		return NewClass(fmt.Sprintf("%v", meta["class-name"]), decoded), nil
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
		class := NewClass(fmt.Sprintf("%v", meta["class-name"]), decoded)
		classes = append(classes, class)
	}
	return &task{
		descr:      descr,
		classes:    classes,
		parameters: params,
	}, nil
}
