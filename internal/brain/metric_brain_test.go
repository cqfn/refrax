package brain

import (
	"testing"

	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
)

func TestMetricBrain_Ask_DelegatesToOrigin(t *testing.T) {
	claim := "Give me good Java code!"
	brain := NewMetricBrain(NewMock(), log.NewMock())
	response, err := brain.Ask(claim)
	assert.NoError(t, err)
	assert.Equal(t, claim, response)
}
