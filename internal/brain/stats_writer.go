package brain

// StatsWriter writes Stats to a prarticular output.
type StatsWriter interface {
	Print(stats *Stats) error
}
