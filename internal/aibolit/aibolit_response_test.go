package aibolit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAibolitResponse_Sanitizes_Response(t *testing.T) {
	expected := []string {
		"x/Foo.java[50]: Non final class (P24: 0.20)",
		"x/Foo.java[476]: Private static method (P25: 7.50)",
		"x/Foo.java[471]: String concat (P17: 1.60)",
	}
	lines := []string {
		"ignore: []",
		"Show pattern with the largest contribution to Cognitive Complexity",
	}
	lines = append(lines, expected...)
	sanitized := NewAibolitResponse(strings.Join(lines, "\n")).Sanitized()
	assert.Equal(t, strings.Join(expected, "\n"), sanitized)
}
