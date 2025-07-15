package brain

import (
	"time"
)

// MetricBrain is a wrapper around a Brain with added functionality
// to track metrics such as the duration of each question asked.
type MetricBrain struct {
	origin Brain
	stats  *Stats
}

// NewMetricBrain creates a new MetricBrain instance wrapping the given Brain
// and using the provided stats for writing statistics.
func NewMetricBrain(brain Brain, s *Stats) Brain {
	return &MetricBrain{brain, s}
}

// Ask sends a question to the underlying Brain and tracks the time
// taken to process the question.
func (b *MetricBrain) Ask(question string) (string, error) {
	start := time.Now()
	result, err := b.origin.Ask(question)
	duration := time.Since(start)
	b.stats.Add(duration)
	return result, err
}
