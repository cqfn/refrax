package client

import (
	"fmt"
	"strings"
)

type Project interface {
	Classes() ([]JavaClass, error)
}

type JavaClass interface {
	Name() string
	Content() string
	SetContent(content string) error
}

type InMemoryProject struct {
	files map[string]JavaClass
}

type InMemoryJavaClass struct {
	name    string
	content string
}

func NewMockProject() Project {
	mapping := map[string]string{
		"Main.java": "public class Main {\n\tpublic static void main(String[] args) {\n\t\tString m = \"Hello, World\";\n\t\tSystem.out.println(m);\n\t}\n}\n",
	}
	return NewInMemoryProject(mapping)
}

func SingleClassProject(name, content string) Project {
	mapping := map[string]string{
		name: content,
	}
	return NewInMemoryProject(mapping)
}

func NewInMemoryProject(files map[string]string) Project {
	res := make(map[string]JavaClass, len(files))
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

func (i *InMemoryProject) Classes() ([]JavaClass, error) {
	res := make([]JavaClass, 0)
	for _, class := range i.files {
		res = append(res, class)
	}
	return res, nil
}

func (i *InMemoryJavaClass) SetContent(content string) error {
	i.content = content
	return nil
}

func (i *InMemoryJavaClass) Content() string {
	return i.content
}

func (i *InMemoryJavaClass) Name() string {
	return i.name
}

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
