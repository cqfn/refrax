package stats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokens_Counts(t *testing.T) {
	tokens, err := Tokens("Hello world!")
	require.NoError(t, err, "Expected to create tokens without error")
	assert.Equal(t, 3, tokens)
}

func TestStats_Add_Success(t *testing.T) {
	s1 := &Stats{
		Name:          "Stats1",
		llmreq:        []time.Duration{time.Millisecond, 2 * time.Millisecond},
		llmreqtokens:  10,
		llmresptokens: 15,
		llmreqbytes:   100,
		llmrespbytes:  200,
		a2areqs:       []time.Duration{3 * time.Millisecond},
		a2areqtokens:  5,
		a2aresptokens: 8,
		a2areqbytes:   50,
		a2arespbytes:  80,
	}
	s2 := &Stats{
		Name:          "Stats2",
		llmreq:        []time.Duration{500 * time.Microsecond},
		llmreqtokens:  20,
		llmresptokens: 25,
		llmreqbytes:   300,
		llmrespbytes:  400,
		a2areqs:       []time.Duration{4 * time.Millisecond},
		a2areqtokens:  15,
		a2aresptokens: 18,
		a2areqbytes:   60,
		a2arespbytes:  90,
	}
	combined := s1.Add(s2)
	require.NotNil(t, combined)
	assert.Equal(t, "Stats1 + Stats2", combined.Name)
	assert.Equal(t, []time.Duration{time.Millisecond, 2 * time.Millisecond, 500 * time.Microsecond}, combined.llmreq)
	assert.Equal(t, []time.Duration{3 * time.Millisecond, 4 * time.Millisecond}, combined.a2areqs)
	assert.Equal(t, 30, combined.llmreqtokens)
	assert.Equal(t, 40, combined.llmresptokens)
	assert.Equal(t, 400, combined.llmreqbytes)
	assert.Equal(t, 600, combined.llmrespbytes)
	assert.Equal(t, 20, combined.a2areqtokens)
	assert.Equal(t, 26, combined.a2aresptokens)
	assert.Equal(t, 110, combined.a2areqbytes)
	assert.Equal(t, 170, combined.a2arespbytes)
}

func TestStats_Add_EmptyOther(t *testing.T) {
	s := &Stats{
		Name:          "Stats1",
		llmreq:        []time.Duration{time.Millisecond},
		llmreqtokens:  10,
		llmresptokens: 15,
		llmreqbytes:   100,
		llmrespbytes:  200,
		a2areqs:       []time.Duration{2 * time.Millisecond},
		a2areqtokens:  5,
		a2aresptokens: 8,
		a2areqbytes:   50,
		a2arespbytes:  80,
	}
	empty := &Stats{}

	combined := s.Add(empty)

	require.NotNil(t, combined)
	assert.Equal(t, "Stats1 + ", combined.Name)
	assert.Equal(t, []time.Duration{time.Millisecond}, combined.llmreq)
	assert.Equal(t, []time.Duration{2 * time.Millisecond}, combined.a2areqs)
	assert.Equal(t, 10, combined.llmreqtokens)
	assert.Equal(t, 15, combined.llmresptokens)
	assert.Equal(t, 100, combined.llmreqbytes)
	assert.Equal(t, 200, combined.llmrespbytes)
	assert.Equal(t, 5, combined.a2areqtokens)
	assert.Equal(t, 8, combined.a2aresptokens)
	assert.Equal(t, 50, combined.a2areqbytes)
	assert.Equal(t, 80, combined.a2arespbytes)
}

func TestStats_Add_EmptySelf(t *testing.T) {
	empty := &Stats{}
	s := &Stats{
		Name:          "Stats2",
		llmreq:        []time.Duration{500 * time.Microsecond},
		llmreqtokens:  20,
		llmresptokens: 25,
		llmreqbytes:   300,
		llmrespbytes:  400,
		a2areqs:       []time.Duration{4 * time.Millisecond},
		a2areqtokens:  15,
		a2aresptokens: 18,
		a2areqbytes:   60,
		a2arespbytes:  90,
	}

	combined := empty.Add(s)

	require.NotNil(t, combined)
	assert.Equal(t, " + Stats2", combined.Name)
	assert.Equal(t, []time.Duration{500 * time.Microsecond}, combined.llmreq)
	assert.Equal(t, []time.Duration{4 * time.Millisecond}, combined.a2areqs)
	assert.Equal(t, 20, combined.llmreqtokens)
	assert.Equal(t, 25, combined.llmresptokens)
	assert.Equal(t, 300, combined.llmreqbytes)
	assert.Equal(t, 400, combined.llmrespbytes)
	assert.Equal(t, 15, combined.a2areqtokens)
	assert.Equal(t, 18, combined.a2aresptokens)
	assert.Equal(t, 60, combined.a2areqbytes)
	assert.Equal(t, 90, combined.a2arespbytes)
}

func TestStats_Add_BothEmpty(t *testing.T) {
	empty1 := &Stats{}
	empty2 := &Stats{}

	combined := empty1.Add(empty2)

	combined.Name = "total"
	require.NotNil(t, combined)
	assert.Equal(t, "total", combined.Name)
	assert.Empty(t, combined.llmreq)
	assert.Empty(t, combined.a2areqs)
	assert.Equal(t, 0, combined.llmreqtokens)
	assert.Equal(t, 0, combined.llmresptokens)
	assert.Equal(t, 0, combined.llmreqbytes)
	assert.Equal(t, 0, combined.llmrespbytes)
	assert.Equal(t, 0, combined.a2areqtokens)
	assert.Equal(t, 0, combined.a2aresptokens)
	assert.Equal(t, 0, combined.a2areqbytes)
	assert.Equal(t, 0, combined.a2arespbytes)
}
