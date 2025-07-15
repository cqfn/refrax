package brain

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
	stats := &Stats{}
	stats.Add(1 * time.Second)
	stats.Add(2 * time.Second)
	stats.Add(3 * time.Second)

	err := w.Print(stats)

	require.NoError(t, err)
	file, err := os.Open(filepath.Clean(p))
	require.NoError(t, err)
	defer func() { _ = file.Close() }()
	lines, err := csv.NewReader(file).ReadAll()
	require.NoError(t, err)
	require.Len(t, lines, 4)
	assert.Equal(t, []string{"Question", "Duration"}, lines[0])
	assert.Equal(t, []string{"1", "1s"}, lines[1])
	assert.Equal(t, []string{"2", "2s"}, lines[2])
	assert.Equal(t, []string{"3", "3s"}, lines[3])
}
