package log

type Prefixed struct {
	prefix   string
	original Logger
}

func NewPrefixed(prefix string, original Logger) Logger {
	return &Prefixed{original: original, prefix: prefix}
}

func (p *Prefixed) Info(msg string, args ...any) {
	p.original.Info(p.prefix+": "+msg, args...)
}

func (p *Prefixed) Debug(msg string, args ...any) {
	p.original.Debug(p.prefix+": "+msg, args...)
}

func (p *Prefixed) Warn(msg string, args ...any) {
	p.original.Warn(p.prefix+": "+msg, args...)
}

func (p *Prefixed) Error(msg string, args ...any) {
	p.original.Error(p.prefix+": "+msg, args...)
}
