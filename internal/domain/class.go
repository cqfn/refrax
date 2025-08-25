package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

// Class represents a code or text entity with a name and content.
type Class interface {
	Name() string
	Content() string
	Path() string
	SetContent(content string) error
}

type memClass struct {
	content string
	path    string
	name    string
}

// NewInMemoryClass creates a new Class instance with the given name and content.
func NewInMemoryClass(name, path, content string) Class {
	return &memClass{
		name:    name,
		path:    path,
		content: content,
	}
}

// Content implements Class.
func (a *memClass) Content() string {
	return a.content
}

// Name implements Class.
func (a *memClass) Name() string {
	return a.name
}

// Path implements Class.
func (a *memClass) Path() string {
	return a.path
}

// SetContent updates the content of the class.
func (a *memClass) SetContent(content string) error {
	a.content = content
	return nil
}

// FSClass represents a Java class file stored in the filesystem.
type FSClass struct {
	name string
	path string
}

// NewFSClass creates a new FilesystemClass with the given name, path, and content.
func NewFSClass(name, path string) *FSClass {
	return &FSClass{
		name: name,
		path: path,
	}
}

func (c *FSClass) Name() string {
	return c.name
}

func (c *FSClass) Path() string {
	return c.path
}

func (c *FSClass) Content() string {
	content, readErr := os.ReadFile(filepath.Clean(c.path))
	if readErr != nil {
		panic(readErr)
	}
	return string(content)
}

// SetContent updates the content of the Java class file and writes it to the filesystem.
func (c *FSClass) SetContent(content string) error {
	err := os.WriteFile(c.path, []byte(content), 0o600)
	if err != nil {
		return fmt.Errorf("error writing content to file %s: %w", c.path, err)
	}
	return nil
}
