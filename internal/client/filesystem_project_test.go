package client

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFilesystemProject(t *testing.T) {
	projectPath := "/example/path"
	project := NewFilesystemProject(projectPath)

	assert.NotNil(t, project)
	assert.Equal(t, projectPath, project.path)
}

func TestFilesystemProject_Classes_Success(t *testing.T) {
	tempDir := t.TempDir()
	first := filepath.Join(tempDir, "Class1.java")
	second := filepath.Join(tempDir, "Class2.java")
	require.NoError(t, os.WriteFile(first, []byte("class Class1 {}"), 0o644))
	require.NoError(t, os.WriteFile(second, []byte("class Class2 {}"), 0o644))
	project := NewFilesystemProject(tempDir)

	classes, err := project.Classes()

	require.NoError(t, err)
	assert.Len(t, classes, 2)

	names := []string{classes[0].Name(), classes[1].Name()}
	assert.Contains(t, names, "Class1")
	assert.Contains(t, names, "Class2")
}

func TestFilesystemProject_Classes_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	project := NewFilesystemProject(tempDir)

	classes, err := project.Classes()

	require.NoError(t, err)
	assert.Empty(t, classes)
}

func TestFilesystemProject_Classes_NonJavaFilesIgnored(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "NotAJavaFile.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("Some content"), 0o644))

	project := NewFilesystemProject(tempDir)

	classes, err := project.Classes()

	require.NoError(t, err)
	assert.Empty(t, classes)
}

func TestFilesystemProject_Classes_ErrorReadingFile(t *testing.T) {
	SkipOnWindows(t)
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "Class1.java")
	require.NoError(t, os.WriteFile(filePath, []byte("class Class1 {}"), 0o644))
	require.NoError(t, os.Chmod(filePath, 0o222)) // Write-only

	project := NewFilesystemProject(tempDir)

	classes, err := project.Classes()

	assert.Nil(t, classes)
	assert.Error(t, err)
}

func TestFilesystemProject_Classes_ErrorDuringTraversal(t *testing.T) {
	SkipOnWindows(t)
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0o000)) // No permissions
	project := NewFilesystemProject(tempDir)

	classes, err := project.Classes()

	assert.Nil(t, classes)
	assert.Error(t, err)
	require.NoError(t, os.Chmod(subDir, 0o755))
}

func TestFilesystemJavaClass_Name(t *testing.T) {
	class := &FilesystemJavaClass{name: "TestClass"}

	assert.Equal(t, "TestClass", class.Name())
}

func TestFilesystemJavaClass_Content(t *testing.T) {
	class := &FilesystemJavaClass{content: "class TestClass {}"}

	assert.Equal(t, "class TestClass {}", class.Content())
}

func TestFilesystemJavaClass_SetContent_Success(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "TestClass.java")
	require.NoError(t, os.WriteFile(filePath, []byte("class TestClass {}"), 0o644))
	class := &FilesystemJavaClass{
		name:    "TestClass",
		content: "class TestClass {}",
		path:    filePath,
	}
	newContent := "class UpdatedClass {}"

	err := class.SetContent(newContent)

	require.NoError(t, err)
	assert.Equal(t, newContent, class.Content())
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, newContent, string(content))
}

func SkipOnWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows")
	}
}
