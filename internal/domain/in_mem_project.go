package domain

import (
	"fmt"
	"strings"
)

// InMemoryProject is an implementation of Project that stores its classes in memory.
type InMemoryProject struct {
	files map[string]Class
}

// NewMock creates a mock project with predefined content for testing purposes.
func NewMock() Project {
	return NewInMemory(NewInMemoryClass("Main.java", ".", "public class Main {\n\tpublic static void main(String[] args) {\n\t\tString m = \"Hello, World\";\n\t\tSystem.out.println(m);\n\t}\n}\n"))
}

// NewInMemory creates a new in-memory project with the given map of file names to Java class content.
func NewInMemory(classes ...Class) Project {
	res := make(map[string]Class, len(classes))
	for _, class := range classes {
		res[class.Path()] = class
	}
	return &InMemoryProject{
		files: res,
	}
}

// Classes retrieves all Java classes in the in-memory project.
func (i *InMemoryProject) Classes() ([]Class, error) {
	res := make([]Class, 0)
	for _, class := range i.files {
		res = append(res, class)
	}
	return res, nil
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
