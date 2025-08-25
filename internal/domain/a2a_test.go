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
			Meta: map[string]any{"key": "value"},
		},
		Classes: []Class{
			NewInMemoryClass("TestClass", "test/path/TestClass.java", "public class TestClass {}"),
			NewInMemoryClass("AnotherClass", "test/path/AnotherClass.java", "public class AnotherClass {}"),
		},
		Examples: []Class{
			NewInMemoryClass("ExampleClass", "test/path/ExampleClass.java", "public class ExampleClass {}"),
		},
		Suggestions: []Suggestion{
			*NewSuggestion("Improve TestClass", "test/path/TestClass.java"),
			*NewSuggestion("Add documentation to AnotherClass", "test/path/AnotherClass.java"),
			*NewSuggestion("Refactor ExampleClass", "test/path/ExampleClass.java"),
		},
	}
	after, err := UnmarshalJob(before.Marshal().Message)
	require.NoError(t, err, "Unmarshaling job should not return an error")
	assert.Equal(t, before, after, "Jobs should be equal after marshaling and unmarshaling")
}
