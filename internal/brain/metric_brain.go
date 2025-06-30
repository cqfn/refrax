package brain

import (
	"time"

	"github.com/cqfn/refrax/internal/log"
)

const questionLength = 64

type MetricBrain struct {
	origin Brain
	stats[] time.Duration
	log log.Logger
}

func NewMetricBrain(brain Brain, log log.Logger) Brain {
	return &MetricBrain{brain, []time.Duration{}, log}
}

func (b *MetricBrain) Ask(question string) (string, error) {
	start := time.Now()
	result, err := b.origin.Ask(question);
	duration := time.Since(start)
	b.stats = append(b.stats, duration)
	return result, err
}

// @todo #21:35min Move `PrintStats` method to the more properly organized abstraction.
//  Currently, we needed it in order to aggregate all the messages in the `stats` map, and then
//  print it. Since there is no `PrintStats` method in original `Brain`, it looks ugly when we call this
//  function on `Brain` instance in `refrax_client.go`. We should organize more proper abstraction
//  around the aggregation and printing of stats.
func (b *MetricBrain) PrintStats() {
	b.log.Info("Total messages asked: %d", len(b.stats))
	for i,d := range b.stats {
		b.log.Info("Brain finished asking question #%d in %s", i+1, d)
	}
}
