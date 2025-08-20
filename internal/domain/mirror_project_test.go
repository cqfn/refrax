package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMirrorProject_Mirrors_Successfully(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	err := os.WriteFile(filepath.Join(srcDir, "Copy.java"), []byte("class Copy{}"), 0o600)
	require.NoError(t, err)

	original := NewFilesystem(srcDir)

	mp, err := NewMirrorProject(original, dstDir)
	require.NoError(t, err)
	require.NotNil(t, mp)

	_, err = os.Stat(filepath.Join(dstDir, "Copy.java"))
	assert.NoError(t, err)
}

func TestNewMirrorProject_CreatesDirectories(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "nonexistent", "dir")

	original := NewFilesystem(srcDir)
	mp, err := NewMirrorProject(original, dstDir)

	require.NoError(t, err)
	require.NotNil(t, mp)
	_, err = os.Stat(dstDir)
	assert.NoError(t, err)
}

func TestMirrorProject_Classes(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	err := os.WriteFile(filepath.Join(src, "App.java"), []byte("public class App{}"), 0o600)
	require.NoError(t, err)

	original := NewFilesystem(src)
	mp, err := NewMirrorProject(original, dst)
	require.NoError(t, err)

	classes, err := mp.Classes()
	assert.NoError(t, err)
	assert.NotNil(t, classes)
}

func TestNewMirrorProject_OverwritesExistingDirectory(t *testing.T) {
	originalDir := t.TempDir()
	originalFile := filepath.Join(originalDir, "Original.java")
	require.NoError(t, os.WriteFile(originalFile, []byte("original content"), 0o600))

	mirrorDir := t.TempDir()
	conflictingFile := filepath.Join(mirrorDir, "Original.java")
	require.NoError(t, os.WriteFile(conflictingFile, []byte("conflicting content"), 0o600))

	orig := NewFilesystem(originalDir)
	_, err := NewMirrorProject(orig, mirrorDir)
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Clean(filepath.Join(mirrorDir, "Original.java")))
	require.NoError(t, err)
	require.Equal(t, "original content", string(data))
}
