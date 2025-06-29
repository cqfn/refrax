package brain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBrainWithStats_Ask_DelegatesToOrigin(t *testing.T) {
	claim := "Give me good Java code!"
	brain := NewBrainWithStats(NewMock(), make(map[string]time.Duration))
	response, err := brain.Ask(claim)
	assert.NoError(t, err)
	assert.Equal(t, "mock response to: " + claim, response)
}
