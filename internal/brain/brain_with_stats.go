package brain

import "github.com/cqfn/refrax/internal/log"

type BrainWithStats struct {
	origin Brain
}

func NewBrainWithStats(brain Brain) Brain {
	return &BrainWithStats{brain}
}

func (b *BrainWithStats) Ask(question string) (string, error) {
	log.Info("Brain with stats asks!")
	return b.origin.Ask(question);
}
