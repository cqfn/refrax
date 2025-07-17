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
	all := stats.Entries()
	for _, v := range all {
		s.log.Info("%s: %s", v.Title, v.Value)
	}
	return nil
}
