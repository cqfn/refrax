package stats

import (
	"github.com/cqfn/refrax/internal/log"
)

type stdWriter struct {
	log log.Logger
}

// NewStdWriter creates an instance of StatsWriter that prints
// statistics to the stdout.
func NewStdWriter(logger log.Logger) Writer {
	return &stdWriter{log: logger}
}

func (s *stdWriter) Print(stats *Stats) error {
	durations := stats.LLMRequests()
	s.log.Info("Total messages asked: %d", len(durations))
	for i, d := range durations {
		s.log.Info("Brain finished asking question #%d in %s", i+1, d)
	}
	return nil
}
