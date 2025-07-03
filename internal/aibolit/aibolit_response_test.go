package aibolit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAibolitResponse_Sanitizes_Response(t *testing.T) {
	expected := "refrax/Foo.java[50]: Non final class (P24: 0.20)"
	sanitized := NewAibolitResponse(
		"ignore: []\nShow pattern with the largest contribution to Cognitive Complexity\nrefrax/Foo.java[50]: Non final class (P24: 0.20)",
	).Sanitized()
	assert.Equal(t, expected, sanitized)
}
