package domain

import (
	"fmt"
	"strconv"

	"github.com/cqfn/refrax/internal/protocol"
	"github.com/google/uuid"
)

// Task represents a unit of work that contains classes and associated parameters.
type Task interface {
	Description() string
	Classes() []Class
	Param(name string) (string, bool)
}

type task struct {
	descr      string
	classes    []Class
	parameters map[string]any
}

// NewTask creates a new Task instance with the given description, classes, and parameters.
func NewTask(description string, classes []Class, parameters map[string]any) Task {
	return &task{
		descr:      description,
		classes:    classes,
		parameters: parameters,
	}
}

func (t *task) Description() string {
	return t.descr
}

func (t *task) Classes() []Class {
	return t.classes
}

func (t *task) Param(name string) (string, bool) {
	if len(t.parameters) == 0 {
		return "", false
	}
	return fmt.Sprintf("%v", t.parameters[name]), true
}

func (t *task) Marshal() *protocol.Message {
	param, ok := t.Param("max-size")
	if !ok {
		param = "200"
	}
	size, err := strconv.Atoi(param)
	if err != nil {
		panic(err)
	}
	msg := protocol.NewMessage().
		WithMessageID(uuid.NewString()).
		AddPart(protocol.NewText(t.Description()).WithMetadata("max-size", size))
	all := t.Classes()
	for _, class := range all {
		name := class.Name()
		path := class.Path()
		msg = msg.AddPart(protocol.NewFileBytes([]byte(class.Content())).WithMetadata("class-name", name).WithMetadata("class-path", path))
	}
	return msg
}
