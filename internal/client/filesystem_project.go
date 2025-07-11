// Package client provides functionality for working with Java source files in a project structure.
package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilesystemProject represents a project stored in the filesystem.
type FilesystemProject struct {
	path string
}

// NewFilesystemProject creates a new FilesystemProject with the given path.
func NewFilesystemProject(path string) *FilesystemProject {
	return &FilesystemProject{path: path}
}

// Classes retrieves all Java classes in the project directory and its subdirectories.
func (p *FilesystemProject) Classes() ([]JavaClass, error) {
	var classes []JavaClass
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".java") {
			content, readErr := os.ReadFile(filepath.Clean(path))
			if readErr != nil {
				return readErr
			}
			classes = append(classes, &FilesystemJavaClass{
				name:    strings.TrimSuffix(info.Name(), ".java"),
				content: string(content),
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

// FilesystemJavaClass represents a Java class file stored in the filesystem.
type FilesystemJavaClass struct {
	name    string
	content string
	path    string
}

// Name returns the name of the Java class.
func (c *FilesystemJavaClass) Name() string {
	return c.name
}

// Content returns the content of the Java class file.
func (c *FilesystemJavaClass) Content() string {
	return c.content
}

// SetContent updates the content of the Java class file and writes it to the filesystem.
func (c *FilesystemJavaClass) SetContent(content string) error {
	c.content = content
	err := os.WriteFile(c.path, []byte(content), 0o600)
	if err != nil {
		return fmt.Errorf("error writing content to file %s: %w", c.path, err)
	}
	return nil
}
