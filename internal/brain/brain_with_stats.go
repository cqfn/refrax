package brain

import (
	"time"

	"github.com/cqfn/refrax/internal/log"
)

type BrainWithStats struct {
	origin Brain
}

func NewBrainWithStats(brain Brain) Brain {
	return &BrainWithStats{brain}
}

func (b *BrainWithStats) Ask(question string) (string, error) {
	log.Info("Brain starts asking...")
	start := time.Now()
	result, err := b.origin.Ask(question);
	log.Info("Brain finished asking in %s", time.Since(start))
	return result, err
}
