package domain

// Class represents a code or text entity with a name and content.
type Class interface {
	Name() string
	Content() string
	Path() string
	SetContent(content string) error
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

// SetContent updates the content of the class.
func (a *class) SetContent(_ string) error {
	panic("SetContent is not implemented for a2a class")
}
