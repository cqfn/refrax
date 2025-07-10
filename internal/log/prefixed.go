package log

type prefixed struct {
	prefix   string
	original Logger
}

// NewPrefixed creates a new Logger that prefixes all log messages with the specified prefix.
func NewPrefixed(prefix string, original Logger) Logger {
	return &prefixed{original: original, prefix: prefix}
}

func (p *prefixed) Info(msg string, args ...any) {
	p.original.Info(p.prefix+": "+msg, args...)
}

func (p *prefixed) Debug(msg string, args ...any) {
	p.original.Debug(p.prefix+": "+msg, args...)
}

func (p *prefixed) Warn(msg string, args ...any) {
	p.original.Warn(p.prefix+": "+msg, args...)
}

func (p *prefixed) Error(msg string, args ...any) {
	p.original.Error(p.prefix+": "+msg, args...)
}
