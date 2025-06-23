package brain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_WithDeepSeekProvider_ReturnsDeepSeekBrain(t *testing.T) {
	provider := "deepseek"
	token := "valid_token"

	result := New(provider, token)

	_, ok := result.(*DeepSeek)
	require.True(t, ok, "Expected result to be of type DeepSeek")
}

func TestNew_WithUnknownProvider_ReturnsMockBrain(t *testing.T) {
	provider := "unknown"
	token := "any_token"

	result := New(provider, token)

	_, ok := result.(*MockBrain)
	require.True(t, ok, "Expected result to be of type Mock")
}
