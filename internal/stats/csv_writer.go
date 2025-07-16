// Package stats provides functionality to write statistics.
package stats

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cqfn/refrax/internal/log"
)

type csvWriter struct {
	path string
}

// NewCSVWriter created an instance of StatsWriter that saves
// statistics in CSV format.
func NewCSVWriter(path string) Writer {
	return &csvWriter{path: path}
}

func (c *csvWriter) Print(stats *Stats) error {
	file, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func() { _ = file.Close() }()
	w := csv.NewWriter(file)
	defer w.Flush()
	if err = w.Write([]string{"Question", "Duration"}); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}
	for i, duration := range stats.LLMRequests() {
		if err = w.Write([]string{strconv.Itoa(i + 1), duration.String()}); err != nil {
			return fmt.Errorf("failed to write row %d: %v", i+1, err)
		}
	}
	abs, err := filepath.Abs(c.path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	log.Info("statistics written to %s", abs)
	return nil
}
