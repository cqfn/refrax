package project

import (
	"fmt"
	"strings"

	"github.com/cqfn/refrax/internal/domain"
)

// InMemoryProject is an implementation of Project that stores its classes in memory.
type InMemoryProject struct {
	files map[string]domain.Class
}

// InMemoryJavaClass is an implementation of JavaClass that stores its data in memory.
type InMemoryJavaClass struct {
	name    string
	content string
}

// NewMock creates a mock project with predefined content for testing purposes.
func NewMock() domain.Project {
	mapping := map[string]string{
		"Main.java": "public class Main {\n\tpublic static void main(String[] args) {\n\t\tString m = \"Hello, World\";\n\t\tSystem.out.println(m);\n\t}\n}\n",
	}
	return NewInMemory(mapping)
}

// SingleClass creates a project containing a single Java class with the provided name and content.
func SingleClass(name, content string) domain.Project {
	mapping := map[string]string{
		name: content,
	}
	return NewInMemory(mapping)
}

// NewInMemory creates a new in-memory project with the given map of file names to Java class content.
func NewInMemory(files map[string]string) domain.Project {
	res := make(map[string]domain.Class, len(files))
	for name, content := range files {
		res[name] = &InMemoryJavaClass{
			name:    name,
			content: content,
		}
	}
	return &InMemoryProject{
		files: res,
	}
}

// Classes retrieves all Java classes in the in-memory project.
func (i *InMemoryProject) Classes() ([]domain.Class, error) {
	res := make([]domain.Class, 0)
	for _, class := range i.files {
		res = append(res, class)
	}
	return res, nil
}

// SetContent updates the content of the in-memory Java class.
func (i *InMemoryJavaClass) SetContent(content string) error {
	i.content = content
	return nil
}

// Content retrieves the content of the in-memory Java class.
func (i *InMemoryJavaClass) Content() string {
	return i.content
}

// Name retrieves the name of the in-memory Java class.
func (i *InMemoryJavaClass) Name() string {
	return i.name
}

// String returns a string representation of the in-memory project.
func (i *InMemoryProject) String() string {
	names := make([]string, 0, len(i.files))
	for name := range i.files {
		names = append(names, name)
	}
	if len(names) == 0 {
		return "[empty project]"
	}
	return fmt.Sprintf("[%s]", strings.Join(names, ", "))
}
