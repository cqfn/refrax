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

// @todo #21:35min Move `PrintStats` method to the more properly organized abstraction.
//  Currently, we needed it in order to aggregate all the messages in the `stats` map, and then
//  print it. Since there is no `PrintStats` method in original `Brain`, it looks ugly when we call this
//  function on `Brain` instance in `refrax_client.go`. We should organize more proper abstraction
//  around the aggregation and printing of stats.
func (b *BrainWithStats) PrintStats() {
	b.log.Info("Total messages asked: %d", len(b.stats))
	for q,d := range b.stats {
		if len(q) > questionLength {
			q = q[:questionLength] + "..."
		}
		b.log.Info("Brain finished asking %q in %s", q, d)
	}
}
