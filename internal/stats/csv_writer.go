// Package stats provides functionality to write statistics.
package stats

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

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

// Print writes the statistics to a CSV file.
func (c *csvWriter) Print(stats *Stats) error {
	file, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func() { _ = file.Close() }()
	w := csv.NewWriter(file)
	defer w.Flush()
	if err = w.Write([]string{"Metric", "Value"}); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}
	for _, v := range stats.Entries() {
		err = w.Write([]string{v.Title, v.Value})
		if err != nil {
			return fmt.Errorf("failed to write entry %s: %v", v.Title, err)
		}
	}
	abs, err := filepath.Abs(c.path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	log.Info("statistics written to %s", abs)
	return nil
}
