package prompts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const system = `# System Prompt for TestAgent

You are a non-chatty tool that works with **Java code**. Respond with the **required output only**, no preface, no explanations, no headings, no extra lines.

## Context
TestProject

## Constraints
- None

## Defaults
- Null verbosity: Do not include rationale or commentary.

## Capabilities
- Basic

## Final Check (silent)
- No explanations, no markdown fences.
`

const user = "# User Prompt for UserAgent\n"

func TestSystem_String_WithValidTemplate(t *testing.T) {
	s := System{
		AgentName:      "TestAgent",
		ProjectContext: "TestProject",
		Constraints:    []string{"None"},
		Capabilities:   []string{"Basic"},
	}

	res := s.String()

	require.NotNil(t, res)
	require.Equal(t, system, strings.ReplaceAll(res, "\r", ""))
}

func TestUser_String_WithValidTemplate(t *testing.T) {
	u := User{
		Data: map[string]any{
			"AgentName": "UserAgent",
		},
		Name: "test.md.tmpl",
	}

	res := u.String()

	require.NotNil(t, res)
	require.Equal(t, user, strings.ReplaceAll(res, "\r", ""))
}

func TestUser_String_FromAnotherFolder(t *testing.T) {
	u := User{
		Data: map[string]any{},
		Name: "critic/critic.md.tmpl",
	}

	res := u.String()

	require.NotNil(t, res)
	require.Contains(t, res, "Analyze the following Java code")
}
