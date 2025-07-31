package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinesAllToolsTogether(t *testing.T) {
	assert.Equal(
		t,
		"foo\nbar",
		NewCombined(NewMock("foo"), NewMock("bar")).Imperfections(),
		"The result of tool combination does not match with expected",
	)
}
