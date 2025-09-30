package log

import "fmt"

type prefixed struct {
	prefix   string
	original Logger
}

var maxPrefix = 0

// NewPrefixed creates a new Logger that prefixes all log messages with the specified prefix.
func NewPrefixed(prefix string, original Logger) Logger {
	maxPrefix = max(len(prefix), maxPrefix)
	return &prefixed{original: original, prefix: prefix}
}

func (p *prefixed) Info(msg string, args ...any) {
	p.original.Info(p.formatPrefix()+msg, args...)
}

func (p *prefixed) Debug(msg string, args ...any) {
	p.original.Debug(p.formatPrefix()+msg, args...)
}

func (p *prefixed) Warn(msg string, args ...any) {
	p.original.Warn(p.formatPrefix()+msg, args...)
}

func (p *prefixed) Error(msg string, args ...any) {
	p.original.Error(p.formatPrefix()+msg, args...)
}

func (p *prefixed) formatPrefix() string {
	return fmt.Sprintf("%-*s", maxPrefix+2, p.prefix+":")
}
