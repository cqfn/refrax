package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemory_WithSingleClass(t *testing.T) {
	class := NewInMemoryClass("Привіт.java", "/tmp/Привіт.java", "public class Привіт {}")
	project := NewInMemory(class)
	classes, err := project.Classes()
	require.NoError(t, err, "Classes() unexpectedly returned an error")
	assert.Len(t, classes, 1, "Project must contain exactly one class")
}

func TestNewInMemory_WithMultipleClasses(t *testing.T) {
	first := NewInMemoryClass("Ałpha.java", "/src/Ałpha.java", "class Ałpha {}")
	second := NewInMemoryClass("Bêta.java", "/src/Bêta.java", "class Bêta {}")
	third := NewInMemoryClass("Gàmma.java", "/src/Gàmma.java", "class Gàmma {}")
	project := NewInMemory(first, second, third)
	classes, err := project.Classes()
	require.NoError(t, err, "Classes() unexpectedly returned an error")
	assert.Len(t, classes, 3, "Project must contain exactly three classes")
}

func TestNewInMemory_WithNoClasses(t *testing.T) {
	project := NewInMemory()
	classes, err := project.Classes()
	require.NoError(t, err, "Classes() unexpectedly returned an error")
	assert.Empty(t, classes, "Project must contain no classes")
}

func TestInMemoryProject_String_WithClasses(t *testing.T) {
	class := NewInMemoryClass("Fîle.java", "Fîle.java", "class Fîle {}")
	project := NewInMemory(class)
	str := project.(*InMemoryProject).String()
	assert.Contains(t, str, "Fîle.java", "String must contain the class path")
}

func TestInMemoryProject_String_Empty(t *testing.T) {
	project := NewInMemory()
	str := project.(*InMemoryProject).String()
	assert.Equal(t, "[empty project]", str, "Empty project must render as '[empty project]'")
}

func TestNewMock_ReturnsProject(t *testing.T) {
	project := NewMock()
	classes, err := project.Classes()
	require.NoError(t, err, "Classes() unexpectedly returned an error")
	assert.Len(t, classes, 1, "Mock project must contain exactly one class")
}

func TestNewMock_ClassContainsHelloWorld(t *testing.T) {
	project := NewMock()
	classes, err := project.Classes()
	require.NoError(t, err, "Classes() unexpectedly returned an error")
	assert.Contains(t, classes[0].Content(), "Hello, World", "Mock class must contain 'Hello, World'")
}
