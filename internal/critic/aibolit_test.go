package critic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRunner struct {
	output string
	err    error
}

func (m *mockRunner) Run(name string, args ...string) ([]byte, error) {
	return []byte(m.output), m.err
}

func TestSanitizedAibolit_Sanitizes_Response(t *testing.T) {
	expected := []string{
		"x/Foo.java[50]: Non final class (P24: 0.20)",
		"x/Foo.java[476]: Private static method (P25: 7.50)",
		"x/Foo.java[471]: String concat (P17: 1.60)",
	}
	lines := []string{
		"ignore: []",
		"Show pattern with the largest contribution to Cognitive Complexity",
	}
	lines = append(lines, expected...)
	aibolit := NewAibolit("Foo.java")
	aibolit.executor = &mockRunner{output: strings.Join(lines, "\n")}

	actual := aibolit.Imperfections()

	assert.Equal(t, strings.Join(expected, "\n"), actual)
}
