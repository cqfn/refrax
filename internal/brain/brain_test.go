package brain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	system  = "system prompt"
	address = "http://custom-provider.com/v1/chat/completions"
)

func TestNew_WithDeepSeekProvider_ReturnsDeepSeekBrain(t *testing.T) {
	token := "valid_token"

	result, err := New(deepseek, address, token, "gemma3", system)

	require.NoError(t, err, "Expected no error when creating DeepSeek brain")
	_, ok := result.(*deepSeek)
	require.True(t, ok, "Expected result to be of type DeepSeek")
}

func TestNew_WithOpenAIProvider_ReturnsOpenAIBrain(t *testing.T) {
	token := "valid_openai_token"

	result, err := New(openai, address, token, "llama", system)

	require.NoError(t, err, "Expected no error when creating OpenAI brain")
	_, ok := result.(*openAI)
	require.True(t, ok, "Expected result to be of type OpenAI")
}

func TestNew_MockProviderNoPlaybook_ReturnsMockInstance(t *testing.T) {
	result, err := New(mock, address, "test-token", "qwen", system)

	require.NoError(t, err)
	_, ok := result.(*mockBrain)
	require.True(t, ok, "Expected result to be of type Mock")
}

func TestNew_UnknownProvider_ReturnsError(t *testing.T) {
	b, err := New("unknown", address, "test-token", "deepseek", system)

	require.Error(t, err)
	assert.Nil(t, b)
	assert.Contains(t, err.Error(), "unknown provider: unknown")
}

func TestNew_CustomProvider_ReturnsCustomInstance(t *testing.T) {
	b, err := New(custom, address, "custom-token", "custom-model", system)

	require.NoError(t, err)
	assert.NotNil(t, b)
	_, ok := b.(*openAI)
	require.True(t, ok, "Expected result to be OpenAI compatible instance")
}

func TestNew_CutomProviderMissingURL_ReturnsError(t *testing.T) {
	b, err := New(custom, "", "custom-token", "custom-model", system)

	require.Error(t, err)
	assert.Nil(t, b)
	require.Contains(t, err.Error(), "custom provider requires a URL")
}

func TestNew_CustomProviderMissingModel_ReturnsError(t *testing.T) {
	b, err := New(custom, address, "custom-token", "", system)

	require.Error(t, err)
	assert.Nil(t, b)
	require.Contains(t, err.Error(), "custom provider requires a model")
}
