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
func (c *csvWriter) Print(stats ...*Stats) error {
	file, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func() { _ = file.Close() }()
	w := csv.NewWriter(file)
	defer w.Flush()
	header := make([]string, 0)
	header = append(header, "metric")
	for _, s := range stats {
		header = append(header, s.Name)
	}
	if err = w.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}
	values := make(map[string][]string, 0)
	order := make([]string, 0)
	for _, s := range stats {
		for _, v := range s.Entries() {
			if _, ok := values[v.Title]; !ok {
				line := make([]string, 0)
				values[v.Title] = line
				order = append(order, v.Title)
			}
			values[v.Title] = append(values[v.Title], v.Value)
		}
	}
	for _, k := range order {
		v := values[k]
		line := make([]string, 0)
		line = append(line, k)
		line = append(line, v...)
		err = w.Write(line)
		if err != nil {
			return fmt.Errorf("failed to write entry %s: %v", k, err)
		}
	}
	abs, err := filepath.Abs(c.path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	log.Info("Statistics written to %s", abs)
	return nil
}
