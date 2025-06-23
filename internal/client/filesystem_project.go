package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemProject struct {
	path string
}

func NewFilesystemProject(path string) *FilesystemProject {
	return &FilesystemProject{path: path}
}

func (p *FilesystemProject) Classes() ([]JavaClass, error) {
	var classes []JavaClass
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".java") {
			content, readErr := os.ReadFile(path)
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

type FilesystemJavaClass struct {
	name    string
	content string
	path    string
}

func (c *FilesystemJavaClass) Name() string {
	return c.name
}

func (c *FilesystemJavaClass) Content() string {
	return c.content
}

func (c *FilesystemJavaClass) SetContent(content string) error {
	c.content = content
	err := os.WriteFile(c.path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing content to file %s: %w", c.path, err)
	}
	return nil
}
