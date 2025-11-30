package brain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCustom(t *testing.T) {
	token := "test_token"
	url := "http://test_url.com/v1/chat/completions"
	model := "test_model"
	system := "test_system"
	custom, err := NewCustom(token, url, model, system)
	require.NoError(t, err)
	assert.NotNil(t, custom)
	expected, err := NewOpenAI(token, url, model, system)
	require.NoError(t, err)
	assert.Equal(t, expected, custom)
}
