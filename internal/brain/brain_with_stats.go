package brain

import (
	"time"

	"github.com/cqfn/refrax/internal/log"
)

const questionLength = 64

type BrainWithStats struct {
	origin Brain
	stats map[string]time.Duration
}

func NewBrainWithStats(brain Brain, stats map[string]time.Duration) Brain {
	return &BrainWithStats{brain, stats}
}

func (b *BrainWithStats) Ask(question string) (string, error) {
	// append each statement into list and then print it in the end
	// message count, for each message calculate its processing time
	// log.Info("Brain starts asking...")
	start := time.Now()
	result, err := b.origin.Ask(question);
	duration := time.Since(start)
	b.stats[question] = duration
	// log.Info("Brain finished asking in %s", duration)
	return result, err
}

func (b *BrainWithStats) PrintStats() {
	log.Info("Total messages asked: %d", len(b.stats))
	for q,d := range b.stats {
		if len(q) > questionLength {
			q = q[:questionLength] + "..."
		}
		log.Info("Question: %q took %s", q, d)
	}
}
