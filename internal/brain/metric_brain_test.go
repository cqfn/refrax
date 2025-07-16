package brain

import (
	"testing"

	"github.com/cqfn/refrax/internal/stats"
	"github.com/stretchr/testify/assert"
)

func TestMetricBrain_Ask_DelegatesToOrigin(t *testing.T) {
	claim := "Give me good Java code!"
	brain := NewMetricBrain(NewMock(), &stats.Stats{})
	response, err := brain.Ask(claim)
	assert.NoError(t, err)
	assert.Equal(t, claim, response)
}
