package brain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAI_Ask_PositiveCase(t *testing.T) {
	server := NewEchoServer(t, "gpt-3.5-turbo", "test_api_key")
	defer server.Close()

	openai, err := NewOpenAIDefault("test_api_key", "openai sys prompt")
	require.NoError(t, err)
	openai.(*openAI).url = server.URL

	answer, err := openai.Ask("This is a test question")

	require.NoError(t, err)
	require.Equal(t, "This is a test question", answer)
}

func TestOpenAI_Ask_NegativeCase(t *testing.T) {
	server := NewErrorServer(t)
	defer server.Close()

	openai, err := NewOpenAIDefault("test_api_key", "openai system prompt")
	require.NoError(t, err)
	openai.(*openAI).url = server.URL

	answer, err := openai.Ask("This is a test question")

	require.Error(t, err)
	require.Empty(t, answer)
}

func TestOpenAI_WrongUrlPrefix(t *testing.T) {
	brain, err := NewOpenAI("token:", "ftp://invalid-url.com", "model", "system")

	assert.Nil(t, brain)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must start with http:// or https://")
}

func TestOpenAI_Ask_EmptyResponse(t *testing.T) {
	brain, err := NewOpenAI("token", "http://invalid-url.com", "model", "system")

	assert.Nil(t, brain)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must end with /v1/chat/completions")
}
