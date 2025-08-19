// Package project provides functionality for working with Java source files in a project structure.
package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqfn/refrax/internal/domain"
)

// FilesystemProject represents a project stored in the filesystem.
type FilesystemProject struct {
	path string
}

// NewFilesystem creates a new FilesystemProject with the given path.
func NewFilesystem(path string) *FilesystemProject {
	return &FilesystemProject{path: path}
}

// Classes retrieves all Java classes in the project directory and its subdirectories.
func (p *FilesystemProject) Classes() ([]domain.Class, error) {
	var classes []domain.Class
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".java") {
			content, readErr := os.ReadFile(filepath.Clean(path))
			if readErr != nil {
				return readErr
			}
			classes = append(classes, &FilesystemClass{
				content: string(content),
				name:    strings.TrimSuffix(info.Name(), ".java"),
				path:    path,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// String returns the string representation of the FilesystemProject.
func (p *FilesystemProject) String() string {
	return fmt.Sprintf("[%s]", p.path)
}

// FilesystemClass represents a Java class file stored in the filesystem.
type FilesystemClass struct {
	name    string
	content string
	path    string
}

// NewFilesystemClass creates a new FilesystemClass with the given name, path, and content.
func NewFilesystemClass(name, path, content string) *FilesystemClass {
	return &FilesystemClass{
		name:    name,
		path:    path,
		content: content,
	}
}

// Name returns the name of the Java class.
func (c *FilesystemClass) Name() string {
	return c.name
}

// Content returns the content of the Java class file.
func (c *FilesystemClass) Content() string {
	return c.content
}

// Path returns the filesystem path of the Java class file.
func (c *FilesystemClass) Path() string {
	return c.path
}

// SetContent updates the content of the Java class file and writes it to the filesystem.
func (c *FilesystemClass) SetContent(content string) error {
	c.content = content
	err := os.WriteFile(c.path, []byte(content), 0o600)
	if err != nil {
		return fmt.Errorf("error writing content to file %s: %w", c.path, err)
	}
	return nil
}
