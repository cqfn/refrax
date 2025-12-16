package reviewer

import (
	"testing"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSuggestions_ParsesValidLines(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := "/path/Tëst.java: add logging\n/path/Ünïcödé.java: remove redundancy"

	suggestions := a.parseSuggestions(output)

	assert.Equal(t, 2, len(suggestions), "must parse all valid lines")
}

func TestParseSuggestions_ExtractsCorrectPath(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := "/path/Clässë.java: fix method"

	suggestions := a.parseSuggestions(output)

	require.Equal(t, 1, len(suggestions), "must parse exactly one suggestion")
	assert.Equal(t, "/path/Clässë.java", suggestions[0].ClassPath, "path must match")
}

func TestParseSuggestions_ExtractsCorrectText(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := "/path/Fïxër.java: add ünïcödé support"

	suggestions := a.parseSuggestions(output)

	require.Equal(t, 1, len(suggestions), "must parse exactly one suggestion")
	assert.Equal(t, "add ünïcödé support", suggestions[0].Text, "text must match")
}

func TestParseSuggestions_SkipsMalformedLines(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := "malformed line without delimiter\n/path/Välïd.java: valid suggestion"

	suggestions := a.parseSuggestions(output)

	assert.Equal(t, 1, len(suggestions), "must skip malformed lines")
}

func TestParseSuggestions_HandlesEmptyOutput(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := ""

	suggestions := a.parseSuggestions(output)

	assert.Equal(t, 0, len(suggestions), "must return empty slice for empty output")
}

func TestParseSuggestions_HandlesMultipleColons(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := "/path/Fïlë.java: fix: method: signature"

	suggestions := a.parseSuggestions(output)

	require.Equal(t, 1, len(suggestions), "must parse exactly one suggestion")
	assert.Equal(t, "fix: method: signature", suggestions[0].Text, "must preserve colons in text")
}

func TestParseSuggestions_TrimsWhitespace(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
	}
	output := "  /path/Späcë.java  :  add logging  "

	suggestions := a.parseSuggestions(output)

	require.Equal(t, 1, len(suggestions), "must parse exactly one suggestion")
	assert.Equal(t, "/path/Späcë.java", suggestions[0].ClassPath, "path must be trimmed")
	assert.Equal(t, "add logging", suggestions[0].Text, "text must be trimmed")
}

func TestReview_ReturnsEmptySuggestionsOnSuccess(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
		cmds:   []string{"true"},
	}

	artifacts, err := a.Review()

	require.NoError(t, err, "unexpected error during review")
	assert.Equal(t, 0, len(artifacts.Suggestions), "must return no suggestions when command succeeds")
}

func TestReview_ReturnsArtifactsWithDescription(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
		cmds:   []string{"true"},
	}

	artifacts, err := a.Review()

	require.NoError(t, err, "unexpected error during review")
	require.NotNil(t, artifacts.Descr, "description must not be nil")
}

func TestReview_HandlesMultipleCommands(t *testing.T) {
	a := &agent{
		logger: log.NewMock(),
		ai:     brain.NewMock(),
		cmds:   []string{"true", "true"},
	}

	artifacts, err := a.Review()

	require.NoError(t, err, "unexpected error during review")
	assert.NotNil(t, artifacts, "artifacts must not be nil")
}
