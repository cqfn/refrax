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

func (w *stdWriter) Print(stats ...*Stats) error {
	for _, s := range stats {
		all := s.Entries()
		for _, v := range all {
			w.log.Info("%s: %s", v.Title, v.Value)
		}
	}
	return nil
}
