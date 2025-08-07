package domain

import (
	"fmt"
)

// Critic represents an interface to review a class and provide suggestions.
type Critic interface {
	Review(class Class) ([]Suggestion, error)
}

// Fixer represents an interface to fix a class based on suggestions and an example.
type Fixer interface {
	Fix(class Class, suggestions []Suggestion, example Class) (Class, error)
}

// Facilitator represents an interface to refactor tasks into multiple classes.
type Facilitator interface {
	Refactor(task Task) ([]Class, error)
}

// Task represents a unit of work that contains classes and associated parameters.
type Task interface {
	Description() string
	Classes() []Class
	Param(name string) (string, bool)
}

// Class represents a code or text entity with a name and content.
type Class interface {
	Name() string
	Content() string
	Path() string
	SetContent(content string) error
}

// Suggestion represents a proposed improvement or fix for a class.
type Suggestion interface {
	Text() string
}

type task struct {
	descr      string
	classes    []Class
	parameters map[string]any
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

// NewTask creates a new Task instance with the given description, classes, and parameters.
func NewTask(description string, classes []Class, parameters map[string]any) Task {
	return &task{
		descr:      description,
		classes:    classes,
		parameters: parameters,
	}
}

type class struct {
	content string
	path    string
	name    string
}

// NewClass creates a new Class instance with the given name and content.
func NewClass(name, path, content string) Class {
	return &class{
		name:    name,
		path:    path,
		content: content,
	}
}

// Content implements Class.
func (a *class) Content() string {
	return a.content
}

// Name implements Class.
func (a *class) Name() string {
	return a.name
}

// Path implements Class.
func (a *class) Path() string {
	return a.path
}

type suggestion struct {
	text string
}

// NewSuggestion creates a new Suggestion instance with the given text.
func NewSuggestion(text string) Suggestion {
	return &suggestion{
		text: text,
	}
}

// Text implements Suggestion.
func (a *suggestion) Text() string {
	return a.text
}

// SetContent updates the content of the class.
func (a *class) SetContent(_ string) error {
	panic("SetContent is not implemented for a2aClass")
}
