package brain

import (
	"fmt"
	"time"

	"github.com/cqfn/refrax/internal/stats"
)

// MetricBrain is a wrapper around a Brain with added functionality
// to track metrics such as the duration of each question asked.
type MetricBrain struct {
	origin Brain
	stats  *stats.Stats
}

// NewMetricBrain creates a new MetricBrain instance wrapping the given Brain
// and using the provided stats for writing statistics.
func NewMetricBrain(brain Brain, s *stats.Stats) Brain {
	return &MetricBrain{brain, s}
}

// Ask sends a question to the underlying Brain and tracks the time
// taken to process the question.
func (b *MetricBrain) Ask(question string) (string, error) {
	start := time.Now()
	result, err := b.origin.Ask(question)
	if err != nil {
		return "", fmt.Errorf("failed to ask question: %w", err)
	}
	duration := time.Since(start)
	reqt, err := stats.Tokens(question)
	if err != nil {
		return "", fmt.Errorf("failed to count tokens for question: %w", err)
	}
	respt, err := stats.Tokens(result)
	if err != nil {
		return "", fmt.Errorf("failed to count tokens for response: %w", err)
	}
	b.stats.LLMReq(duration, reqt, respt, len(question), len(result))
	return result, err
}
