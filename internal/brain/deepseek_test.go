// brain/deepseek_test.go
package brain

import (
	"testing"
	
	"github.com/stretchr/testify/require"
)

func TestDeepSeek_Ask_PositiveCase(t *testing.T) {
	server := NewEchoServer(t, "deepseek-chat", "test_api_key")
	defer server.Close()

	deepseek := NewDeepSeek("test_api_key")
	deepseek.(*DeepSeek).url = server.URL

	answer, err := deepseek.Ask("This is a test question")
	require.NoError(t, err)
	require.Equal(t, "This is a test question", answer)
}

func TestDeepSeek_Ask_NegativeCase(t *testing.T) {
	server := NewErrorServer(t)
	defer server.Close()
	
	deepseek := NewDeepSeek("test_api_key")
	deepseek.(*DeepSeek).url = server.URL

	answer, err := deepseek.Ask("This is a test question")

	require.Error(t, err)
	require.Empty(t, answer)
}
