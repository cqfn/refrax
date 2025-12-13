package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask_StoresDescription(t *testing.T) {
	task := NewTask("Rèfáctør clàss", nil, nil)
	assert.Equal(t, "Rèfáctør clàss", task.Description(), "Task description must match the provided value")
}

func TestNewTask_StoresClasses(t *testing.T) {
	classes := []Class{
		NewInMemoryClass("Önë.java", "/Önë.java", "class Önë {}"),
		NewInMemoryClass("Twö.java", "/Twö.java", "class Twö {}"),
	}
	task := NewTask("description", classes, nil)
	assert.Len(t, task.Classes(), 2, "Task must contain exactly two classes")
}

func TestNewTask_EmptyClasses(t *testing.T) {
	task := NewTask("empty task", []Class{}, nil)
	assert.Empty(t, task.Classes(), "Task must contain no classes")
}

func TestTask_Param_ReturnsValue(t *testing.T) {
	params := map[string]any{"kéy": "vàlüë"}
	task := NewTask("task", nil, params)
	value, ok := task.Param("kéy")
	assert.True(t, ok, "Param must return true when key exists")
	assert.Equal(t, "vàlüë", value, "Param value must match the stored value")
}

func TestTask_Param_ReturnsEmptyOnMissingKey(t *testing.T) {
	params := map[string]any{"existing": "value"}
	task := NewTask("task", nil, params)
	value, ok := task.Param("nönëxïstënt")
	assert.True(t, ok, "Param returns true even for missing keys when params exist")
	assert.Equal(t, "<nil>", value, "Param returns '<nil>' for missing keys")
}

func TestTask_Param_ReturnsFalseOnNilParams(t *testing.T) {
	task := NewTask("task", nil, nil)
	_, ok := task.Param("àny")
	assert.False(t, ok, "Param must return false when parameters are nil")
}

func TestTask_Param_NumericValue(t *testing.T) {
	params := map[string]any{"nümbër": 42}
	task := NewTask("task", nil, params)
	value, ok := task.Param("nümbër")
	assert.True(t, ok, "Param must return true for numeric parameter")
	assert.Equal(t, "42", value, "Numeric param must be converted to string")
}

func TestTask_Marshal_ContainsDescription(t *testing.T) {
	task := NewTask("Тест описания", nil, nil).(*task)
	msg := task.Marshal()
	assert.NotNil(t, msg, "Marshal must return non-nil message")
}

func TestTask_Marshal_ContainsClasses(t *testing.T) {
	classes := []Class{
		NewInMemoryClass("Clàss.java", "/Clàss.java", "class Clàss {}"),
	}
	task := NewTask("task", classes, nil).(*task)
	msg := task.Marshal()
	assert.NotNil(t, msg, "Marshal must return non-nil message")
}

func TestTask_Marshal_UsesMaxSizeParam(t *testing.T) {
	params := map[string]any{"max-size": "100"}
	task := NewTask("task", nil, params).(*task)
	msg := task.Marshal()
	assert.NotNil(t, msg, "Marshal must return non-nil message with max-size param")
}
