package brain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsk_ValidQuestion(t *testing.T) {
	brain := NewMock()
	question := "What is the capital of France?"

	response, err := brain.Ask(question)

	require.NoError(t, err)
	require.Equal(t, "mock response to: What is the capital of France?", response)
}

func TestAsk_EmptyQuestion(t *testing.T) {
	brain := NewMock()
	question := ""

	response, err := brain.Ask(question)

	require.Error(t, err)
	require.Equal(t, "question cannot be empty", err.Error())
	require.Empty(t, response)
}
