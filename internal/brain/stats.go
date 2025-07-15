package brain

import "time"

// Stats is a struct that contain interaction statistics between components.
type Stats struct {
	durations []time.Duration
}

// Add request timing.
func (s *Stats) Add(duration time.Duration) {
	s.durations = append(s.durations, duration)
}

// Durations retrieved all request timings.
func (s *Stats) Durations() []time.Duration {
	duplicate := make([]time.Duration, len(s.durations))
	copy(duplicate, s.durations)
	return duplicate
}
