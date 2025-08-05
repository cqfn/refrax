package stats

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCSVWriter(t *testing.T) {
	require.NotNil(t, NewCSVWriter(filepath.Join(t.TempDir(), "test_path.csv")))
}

func TestCSVWriter_Print_Success(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "output.csv")
	w := NewCSVWriter(p)
	stats := Stats{Name: "test-stats"}
	stats.LLMReq(1*time.Second, 0, 0, 0, 0)
	stats.LLMReq(2*time.Second, 0, 0, 0, 0)
	stats.LLMReq(3*time.Second, 0, 0, 0, 0)

	err := w.Print(&stats)

	require.NoError(t, err)
	file, err := os.Open(filepath.Clean(p))
	require.NoError(t, err)
	defer func() { _ = file.Close() }()
	lines, err := csv.NewReader(file).ReadAll()
	require.NoError(t, err)
	require.Len(t, lines, 27)
	assert.Equal(t, []string{"metric", "test-stats"}, lines[0])
	assert.Equal(t, []string{"Total LLM messages asked", "3"}, lines[1])
	assert.Equal(t, []string{"Total LLM request duration", "6s"}, lines[2])
	assert.Equal(t, []string{"Total LLM tokens", "0"}, lines[3])
}

func TestCSVWriter_PrintSeveral_Success(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "output.csv")
	w := NewCSVWriter(p)
	first := Stats{Name: "first-stats"}
	first.LLMReq(3*time.Second, 0, 0, 0, 0)
	second := Stats{Name: "second-stats"}
	second.LLMReq(3*time.Second, 1, 1, 1, 1)

	err := w.Print(&first, &second)

	require.NoError(t, err)
	file, err := os.Open(filepath.Clean(p))
	require.NoError(t, err)
	defer func() { _ = file.Close() }()
	lines, err := csv.NewReader(file).ReadAll()
	require.NoError(t, err)
	require.Len(t, lines, 27)
	assert.Equal(t, []string{"metric", "first-stats", "second-stats"}, lines[0])
	assert.Equal(t, []string{"Total LLM messages asked", "1", "1"}, lines[1])
	assert.Equal(t, []string{"Total LLM request duration", "3s", "3s"}, lines[2])
	assert.Equal(t, []string{"Total LLM tokens", "0", "2"}, lines[3])
}
