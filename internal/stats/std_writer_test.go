package stats

import (
	"testing"
	"time"

	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStdWriterCreatesInstance(t *testing.T) {
	w := NewStdWriter(log.NewMock())
	require.NotNil(t, w)
	assert.IsType(t, &stdWriter{}, w)
}

func TestStdWriterPrintLogsCorrectMessageCounts(t *testing.T) {
	m := log.NewMock()
	w := NewStdWriter(m)
	var stats Stats
	stats.LLMReq(1*time.Millisecond, 0, 0, 0, 0)
	stats.LLMReq(2*time.Millisecond, 0, 0, 0, 0)
	stats.LLMReq(3*time.Millisecond, 0, 0, 0, 0)

	err := w.Print(&stats)

	require.NoError(t, err)
	entries := m.Messages
	assert.Equal(t, "mock info: Total LLM messages asked: 3", entries[0])
	assert.Equal(t, "mock info: Total LLM request duration: 6ms", entries[1])
	assert.Equal(t, "mock info: Total LLM tokens: 0", entries[2])
	assert.Equal(t, "mock info: Total LLM request tokens: 0", entries[3])
}

func TestStdWriterPrintHandlesEmptyStats(t *testing.T) {
	m := log.NewMock()
	w := NewStdWriter(m)
	var stats Stats

	err := w.Print(&stats)

	require.NoError(t, err)
	entries := m.Messages
	require.Len(t, entries, 26)
	assert.Equal(t, "mock info: Total LLM messages asked: 0", entries[0])
}

func TestStdWriterPrintLogsCorrectDurations(t *testing.T) {
	m := log.NewMock()
	w := NewStdWriter(m)
	var stats Stats
	stats.LLMReq(1*time.Millisecond, 0, 0, 0, 0)
	stats.LLMReq(1*time.Second, 0, 0, 0, 0)

	err := w.Print(&stats)

	require.NoError(t, err)
	entries := m.Messages
	assert.Contains(t, entries[1], "mock info: Total LLM request duration")
	assert.Contains(t, entries[2], "mock info: Total LLM tokens")
}
