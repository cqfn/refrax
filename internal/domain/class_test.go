package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryClass_StoresName(t *testing.T) {
	class := NewInMemoryClass("Ünîcödé.java", "/path/Ünîcödé.java", "content")
	assert.Equal(t, "Ünîcödé.java", class.Name(), "Class name must match the provided value")
}

func TestNewInMemoryClass_StoresPath(t *testing.T) {
	class := NewInMemoryClass("Tëst.java", "/ûsr/src/Tëst.java", "content")
	assert.Equal(t, "/ûsr/src/Tëst.java", class.Path(), "Class path must match the provided value")
}

func TestNewInMemoryClass_StoresContent(t *testing.T) {
	content := "public class Mÿ { String ñ = \"Héllo\"; }"
	class := NewInMemoryClass("Mÿ.java", "/src/Mÿ.java", content)
	assert.Equal(t, content, class.Content(), "Class content must match the provided value")
}

func TestMemClass_SetContent_UpdatesContent(t *testing.T) {
	class := NewInMemoryClass("Çlàss.java", "/Çlàss.java", "old content")
	err := class.SetContent("new çöntênt")
	require.NoError(t, err, "SetContent unexpectedly returned an error")
	assert.Equal(t, "new çöntênt", class.Content(), "Content must reflect the updated value")
}

func TestMemClass_SetContent_ReturnsNil(t *testing.T) {
	class := NewInMemoryClass("Àny.java", "/Àny.java", "initial")
	err := class.SetContent("updated")
	assert.NoError(t, err, "SetContent must not return an error for in-memory class")
}

func TestNewFSClass_StoresName(t *testing.T) {
	class := NewFSClass("Fïlé.java", "/tmp/Fïlé.java")
	assert.Equal(t, "Fïlé.java", class.Name(), "Class name must match the provided value")
}

func TestNewFSClass_StoresPath(t *testing.T) {
	class := NewFSClass("Päth.java", "/vàr/Päth.java")
	assert.Equal(t, "/vàr/Päth.java", class.Path(), "Class path must match the provided value")
}

func TestFSClass_Content_ReadsFromFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Rëàd.java")
	expected := "class Rëàd { String µ = \"tëst\"; }"
	require.NoError(t, os.WriteFile(path, []byte(expected), 0o600), "failed to write test file")
	class := NewFSClass("Rëàd.java", path)
	assert.Equal(t, expected, class.Content(), "Content must match the file contents")
}

func TestFSClass_SetContent_WritesToFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "Wrïtë.java")
	require.NoError(t, os.WriteFile(path, []byte("initial"), 0o600), "failed to write test file")
	class := NewFSClass("Wrïtë.java", path)
	expected := "class Wrïtë { int ä = 42; }"
	err := class.SetContent(expected)
	require.NoError(t, err, "SetContent unexpectedly returned an error")
	content, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read test file")
	assert.Equal(t, expected, string(content), "File content must match the updated value")
}

func TestFSClass_SetContent_ErrorOnInvalidPath(t *testing.T) {
	class := NewFSClass("Ïnvàlîd.java", "/nönëxîstënt/Ïnvàlîd.java")
	err := class.SetContent("content")
	assert.Error(t, err, "SetContent must return an error for invalid path")
}
