package brain

import (
	"testing"
	"time"

	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
)

func TestBrainWithStats_Ask_DelegatesToOrigin(t *testing.T) {
	claim := "Give me good Java code!"
	brain := NewBrainWithStats(NewMock(), make(map[string]time.Duration), log.NewMock())
	response, err := brain.Ask(claim)
	assert.NoError(t, err)
	assert.Equal(t, "mock response to: " + claim, response)
}
