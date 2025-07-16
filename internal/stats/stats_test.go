package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokens_Counts(t *testing.T) {
	tokens, err := Tokens("Hello world!")
	require.NoError(t, err, "Expected to create tokens without error")
	assert.Equal(t, 3, tokens)
}
