package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FSProj represents a project stored in the filesystem.
type FSProj struct {
	path string
}

// NewFilesystem creates a new FilesystemProject with the given path.
func NewFilesystem(path string) *FSProj {
	return &FSProj{path: path}
}

// Classes retrieves all Java classes in the project directory and its subdirectories.
func (p *FSProj) Classes() ([]Class, error) {
	var classes []Class
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := ".java"
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			c := NewFSClass(strings.TrimSuffix(info.Name(), ext), path)
			classes = append(classes, c)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// String returns the string representation of the FilesystemProject.
func (p *FSProj) String() string {
	return fmt.Sprintf("[%s]", p.path)
}
