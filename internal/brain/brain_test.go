package brain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_WithDeepSeekProvider_ReturnsDeepSeekBrain(t *testing.T) {
	token := "valid_token"

	result := New(deepseek, token)

	_, ok := result.(*deepSeek)
	require.True(t, ok, "Expected result to be of type DeepSeek")
}

func TestNew_WithUnknownProvider_ReturnsMockBrain(t *testing.T) {
	provider := "unknown"
	token := "any_token"

	result := New(provider, token)

	_, ok := result.(*MockBrain)
	require.True(t, ok, "Expected result to be of type Mock")
}

func TestNew_WithOpenAIProvider_ReturnsOpenAIBrain(t *testing.T) {
	token := "valid_openai_token"

	result := New(openai, token)

	_, ok := result.(*openAI)
	require.True(t, ok, "Expected result to be of type OpenAI")
}
