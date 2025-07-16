package stats

// Writer writes Stats to a prarticular output.
type Writer interface {
	Print(stats *Stats) error
}
