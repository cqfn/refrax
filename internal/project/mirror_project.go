package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cqfn/refrax/internal/domain"
)

// MirrorProject decorates FilesystemProject with a mirror location to avoid modifying the original.
type MirrorProject struct {
	mirror Project
}

// NewMirrorProject creates a mirror of the original FilesystemProject at the given path.
func NewMirrorProject(original *FilesystemProject, mirrorPath string) (*MirrorProject, error) {
	if err := os.RemoveAll(filepath.Clean(mirrorPath)); err != nil {
		return nil, fmt.Errorf("failed to remove existing mirror path: %w", err)
	}
	err := os.CopyFS(mirrorPath, os.DirFS(original.path))
	if err != nil {
		return nil, fmt.Errorf("failed to copy project: %w", err)
	}
	return &MirrorProject{mirror: NewFilesystem(mirrorPath)}, nil
}

// Classes retrieves all Java classes from the mirrored project.
func (m *MirrorProject) Classes() ([]domain.Class, error) {
	return m.mirror.Classes()
}
