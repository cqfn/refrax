package brain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenAI_Ask_PositiveCase(t *testing.T) {
	server := NewEchoServer(t, "gpt-3.5-turbo", "test_api_key")
	defer server.Close()

	openai := NewOpenAI("test_api_key")
	openai.(*openAI).url = server.URL

	answer, err := openai.Ask("This is a test question")

	require.NoError(t, err)
	require.Equal(t, "This is a test question", answer)
}

func TestOpenAI_Ask_NegativeCase(t *testing.T) {
	server := NewErrorServer(t)
	defer server.Close()

	openai := NewOpenAI("test_api_key")
	openai.(*openAI).url = server.URL

	answer, err := openai.Ask("This is a test question")

	require.Error(t, err)
	require.Empty(t, answer)
}
