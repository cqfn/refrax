package brain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEchoPlaybook_ReturnsInstance(t *testing.T) {
	ep := NewEchoPlaybook()

	require.NotNil(t, ep, "expected NewEchoPlaybook to return a non-nil instance")
	assert.IsType(t, &echoPlaybook{}, ep, "expected NewEchoPlaybook to return a pointer to echoPlaybook")
}

func TestEchoPlaybook_Ask_ReturnsSameString(t *testing.T) {
	ep := NewEchoPlaybook()
	q := "test question"

	r := ep.Ask(q)

	require.NotNil(t, r, "expected Ask to return a non-nil string")
	assert.Equal(t, q, r, "expected Ask to return the same string as the input")
}

func TestEchoPlaybook_Ask_HandlesEmptyString(t *testing.T) {
	ep := NewEchoPlaybook()
	q := ""

	r := ep.Ask(q)

	require.NotNil(t, r, "expected Ask to return a non-nil string for empty input")
	assert.Equal(t, q, r, "expected Ask to return the same string for empty input")
}

func TestEchoPlaybook_Ask_HandlesSpecialCharacters(t *testing.T) {
	ep := NewEchoPlaybook()
	q := "?!@#$%^&*()"

	r := ep.Ask(q)

	require.NotNil(t, r, "expected Ask to return a non-nil string for special characters")
	assert.Equal(t, q, r, "expected Ask to return the same string for special characters")
}

func TestEchoPlaybook_Ask_HandlesWhitespace(t *testing.T) {
	ep := NewEchoPlaybook()
	q := "   "

	r := ep.Ask(q)

	require.NotNil(t, r, "expected Ask to return a non-nil string for whitespace")
	assert.Equal(t, q, r, "expected Ask to return the same string for whitespace")
}
