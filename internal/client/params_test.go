package client

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockParams_CreatesValidInstance(t *testing.T) {
	params := NewMockParams()
	require.NotNil(t, params, "params must not be nil")
}

func TestNewMockParams_SetsMockProvider(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "mock", params.Provider, "provider must be mock")
}

func TestNewMockParams_SetsDefaultToken(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "ABC", params.Token, "token must be ABC")
}

func TestNewMockParams_EnablesMockProject(t *testing.T) {
	params := NewMockParams()
	assert.True(t, params.MockProject, "mock project must be enabled")
}

func TestNewMockParams_DisablesDebug(t *testing.T) {
	params := NewMockParams()
	assert.False(t, params.Debug, "debug must be disabled")
}

func TestNewMockParams_DisablesStats(t *testing.T) {
	params := NewMockParams()
	assert.False(t, params.Stats, "stats must be disabled")
}

func TestNewMockParams_SetsStdFormat(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "std", params.Format, "format must be std")
}

func TestNewMockParams_SetsDefaultMaxSize(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, 200, params.MaxSize, "max size must be 200")
}

func TestNewMockParams_SetsDiscardLog(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, io.Discard, params.Log, "log must be io.Discard")
}

func TestNewMockParams_SetsDefaultChecks(t *testing.T) {
	params := NewMockParams()
	require.Equal(t, 1, len(params.Checks), "must have one check")
	assert.Equal(t, "mvn clean test", params.Checks[0], "check must be mvn clean test")
}

func TestNewMockParams_DisablesColorless(t *testing.T) {
	params := NewMockParams()
	assert.False(t, params.Colorless, "colorless must be disabled")
}

func TestNewMockParams_SetsDefaultModel(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "gpt-3.5-turbo", params.Model, "model must be gpt-3.5-turbo")
}

func TestNewMockParams_SetsDefaultAttempts(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, 3, params.Attempts, "attempts must be 3")
}

func TestNewMockParams_SetsEmptyPlaybook(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "", params.Playbook, "playbook must be empty")
}

func TestNewMockParams_SetsStatsOutput(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "stats", params.Soutput, "soutput must be stats")
}

func TestNewMockParams_SetsEmptyInput(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "", params.Input, "input must be empty")
}

func TestNewMockParams_SetsEmptyOutput(t *testing.T) {
	params := NewMockParams()
	assert.Equal(t, "", params.Output, "output must be empty")
}

func TestMask_ReturnsEmptyForEmptyToken(t *testing.T) {
	result := mask("")
	assert.Equal(t, "", result, "mask of empty string must be empty")
}

func TestMask_MasksSingleCharacterToken(t *testing.T) {
	result := mask("A")
	assert.Equal(t, "A", result, "single character token must remain fully visible")
}

func TestMask_MasksShortToken(t *testing.T) {
	result := mask("AB")
	assert.Equal(t, "AB", result, "token shorter than 3 characters must remain fully visible")
}

func TestMask_MasksThreeCharacterToken(t *testing.T) {
	result := mask("ABC")
	assert.Equal(t, "ABC", result, "three character token must remain fully visible")
}

func TestMask_MasksFourCharacterToken(t *testing.T) {
	result := mask("ABCD")
	assert.Equal(t, "ABC*", result, "four character token must show first three characters and mask one")
}

func TestMask_MasksLongToken(t *testing.T) {
	result := mask("ABCDEFGH")
	assert.Equal(t, "ABC*****", result, "long token must be masked")
}

func TestMask_ShowsThreeCharacters(t *testing.T) {
	result := mask("secrettoken12345")
	assert.NotEmpty(t, result, "masked token must not be empty")
	assert.Equal(t, "sec*************", result, "token must show first three characters")
}
