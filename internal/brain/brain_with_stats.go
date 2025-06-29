package brain

import (
	"time"

	"github.com/cqfn/refrax/internal/log"
)

const questionLength = 64

type BrainWithStats struct {
	origin Brain
	stats map[string]time.Duration
	log log.Logger
}

func NewBrainWithStats(brain Brain, stats map[string]time.Duration, log log.Logger) Brain {
	return &BrainWithStats{brain, stats, log}
}

func (b *BrainWithStats) Ask(question string) (string, error) {
	start := time.Now()
	result, err := b.origin.Ask(question);
	duration := time.Since(start)
	b.stats[question] = duration
	return result, err
}

func (b *BrainWithStats) PrintStats() {
	b.log.Info("Total messages asked: %d", len(b.stats))
	for q,d := range b.stats {
		if len(q) > questionLength {
			q = q[:questionLength] + "..."
		}
		b.log.Info("Brain finished asking %q in %s", q, d)
	}
}
