package brain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_WithDeepSeekProvider_ReturnsDeepSeekBrain(t *testing.T) {
	token := "valid_token"

	result, err := New(deepseek, token)

	require.NoError(t, err, "Expected no error when creating DeepSeek brain")
	_, ok := result.(*deepSeek)
	require.True(t, ok, "Expected result to be of type DeepSeek")
}

func TestNew_WithOpenAIProvider_ReturnsOpenAIBrain(t *testing.T) {
	token := "valid_openai_token"

	result, err := New(openai, token)

	require.NoError(t, err, "Expected no error when creating OpenAI brain")
	_, ok := result.(*openAI)
	require.True(t, ok, "Expected result to be of type OpenAI")
}

func TestNew_MockProviderNoPlaybook_ReturnsMockInstance(t *testing.T) {
	result, err := New(mock, "test-token")

	require.NoError(t, err)
	_, ok := result.(*mockBrain)
	require.True(t, ok, "Expected result to be of type Mock")
}

func TestNew_UnknownProvider_ReturnsError(t *testing.T) {
	b, err := New("unknown", "test-token")

	require.Error(t, err)
	assert.Nil(t, b)
	assert.Contains(t, err.Error(), "unknown provider: unknown")
}
