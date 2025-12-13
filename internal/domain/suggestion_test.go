package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOldSuggestion_StoresText(t *testing.T) {
	suggestion := NewOldSuggestion("Rèmøvé unused impört")
	assert.Equal(t, "Rèmøvé unused impört", suggestion.Text(), "Suggestion text must match the provided value")
}

func TestNewOldSuggestion_PreservesUnicode(t *testing.T) {
	text := "Добавить метод для обработки данных"
	suggestion := NewOldSuggestion(text)
	assert.Equal(t, text, suggestion.Text(), "Suggestion must preserve unicode characters")
}

func TestNewOldSuggestion_EmptyText(t *testing.T) {
	suggestion := NewOldSuggestion("")
	assert.Empty(t, suggestion.Text(), "Suggestion text must be empty when created with empty string")
}

func TestNewSuggestion_StoresText(t *testing.T) {
	suggestion := NewSuggestion("Ådd null chëck", "/src/Mäïn.java")
	assert.Equal(t, "Ådd null chëck", suggestion.Text, "Suggestion text must match the provided value")
}

func TestNewSuggestion_StoresPath(t *testing.T) {
	suggestion := NewSuggestion("Rêfàctör mëthöd", "/pröjëct/Clàss.java")
	assert.Equal(t, "/pröjëct/Clàss.java", suggestion.ClassPath, "Suggestion path must match the provided value")
}

func TestNewSuggestion_PreservesMultilineText(t *testing.T) {
	text := "Línë önë\nLínë twö\nLínë thrëë"
	suggestion := NewSuggestion(text, "/path/Fïlë.java")
	assert.Equal(t, text, suggestion.Text, "Suggestion must preserve multiline text")
}
