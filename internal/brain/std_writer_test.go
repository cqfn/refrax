package brain

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
	stats := &Stats{}
	stats.Add(1 * time.Millisecond)
	stats.Add(2 * time.Millisecond)
	stats.Add(3 * time.Millisecond)

	err := w.Print(stats)

	require.NoError(t, err)
	entries := m.Messages
	assert.Equal(t, "mock info: Total messages asked: 3", entries[0])
	assert.Equal(t, "mock info: Brain finished asking question #1 in 1ms", entries[1])
	assert.Equal(t, "mock info: Brain finished asking question #2 in 2ms", entries[2])
	assert.Equal(t, "mock info: Brain finished asking question #3 in 3ms", entries[3])
}

func TestStdWriterPrintHandlesEmptyStats(t *testing.T) {
	m := log.NewMock()
	w := NewStdWriter(m)
	stats := &Stats{}

	err := w.Print(stats)

	require.NoError(t, err)
	entries := m.Messages
	require.Len(t, entries, 1)
	assert.Equal(t, "mock info: Total messages asked: 0", entries[0])
}

func TestStdWriterPrintLogsCorrectDurations(t *testing.T) {
	m := log.NewMock()
	w := NewStdWriter(m)
	stats := &Stats{}
	stats.Add(1 * time.Millisecond)
	stats.Add(1 * time.Second)

	err := w.Print(stats)

	require.NoError(t, err)
	entries := m.Messages
	assert.Contains(t, entries[1], "mock info: Brain finished asking question #1 in 1ms")
	assert.Contains(t, entries[2], "mock info: Brain finished asking question #2 in 1s")
}
