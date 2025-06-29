package brain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrainWithStats_Ask_DelegatesToOrigin(t *testing.T) {
	claim := "Give me good Java code!"
	brain := NewBrainWithStats(NewMock())
	response, err := brain.Ask(claim)
	assert.NoError(t, err)
	assert.Equal(t, "mock response to: " + claim, response)
}
