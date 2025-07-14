package facilitator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const clazz = "MyClass"

func TestClassNameHandlesPrefixAndExtractsClassName(t *testing.T) {
	task := "Refactor the class '" + clazz + "'"
	expected := clazz
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected class name to match extracted value")
}

func TestClassNameReturnsEmptyStringWhenNoQuotes(t *testing.T) {
	task := "Refactor the class " + clazz
	expected := ""
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected empty string when no quotes are present")
}

func TestClassNameHandlesEmptyTask(t *testing.T) {
	task := ""
	expected := ""
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected empty string for empty task input")
}

func TestClassNameHandlesMultipleQuotes(t *testing.T) {
	task := "Refactor the class '" + clazz + "' and 'AnotherClass'"
	expected := clazz
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected all quoted text including multiple quotes to be extracted")
}

func TestClassNameHandlesNoPrefix(t *testing.T) {
	task := "Fix the method in class '" + clazz + "'"
	expected := clazz
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected quoted class name to be extracted even if prefix is missing")
}

func TestClassNameHandlesNestedQuotes(t *testing.T) {
	task := "Refactor the class '" + clazz + "'SomeText'AnotherClass'"
	expected := clazz
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected nested quotes to be handled correctly")
}

func TestClassNameHandlesTrimmedPrefix(t *testing.T) {
	task := "Refactor the class ''"
	expected := ""
	actual := className(task)
	assert.Equal(t, expected, actual, "Expected empty string for empty quoted class name after prefix removal")
}

func TestClassNameHandlesTemporaryDirectoryInTask(t *testing.T) {
	dir := t.TempDir()
	task := "Refactor the class '" + dir + "'"
	expected := dir
	actual := className(task)
	require.Equal(t, expected, actual, "Expected directory name in quotes to be extracted correctly")
}
