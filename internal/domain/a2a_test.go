package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalAndUnmarshalJob(t *testing.T) {
	before := &Job{
		Descr: &Description{
			Text: "Test Job",
			meta: map[string]any{"key": "value"},
		},
		Classes: []Class{
			NewClass("TestClass", "test/path/TestClass.java", "public class TestClass {}"),
			NewClass("AnotherClass", "test/path/AnotherClass.java", "public class AnotherClass {}"),
		},
		Examples: []Class{
			NewClass("ExampleClass", "test/path/ExampleClass.java", "public class ExampleClass {}"),
		},
		Suggestions: []Suggestion{
			NewSuggestion("Improve TestClass"),
			NewSuggestion("Add documentation to AnotherClass"),
			NewSuggestion("Refactor ExampleClass"),
		},
	}
	after, err := UnmarshalJob(before.Marshal().Message)
	require.NoError(t, err, "Unmarshaling job should not return an error")
	assert.Equal(t, before, after, "Jobs should be equal after marshaling and unmarshaling")
}
