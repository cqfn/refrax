package log

import (
	"fmt"
)

type colored struct {
	base  Logger
	color Color
}

// Color represents a terminal color escape code.
type Color string

// Predefined terminal color escape codes.
const (
	Black   Color = "\033[30m"       // Black color
	Red     Color = "\033[31m"       // Red color
	Green   Color = "\033[32m"       // Green color
	Yellow  Color = "\033[33m"       // Yellow color
	Blue    Color = "\033[34m"       // Blue color
	Magenta Color = "\033[35m"       // Magenta color
	Cyan    Color = "\033[36m"       // Cyan color
	White   Color = "\033[37m"       // White color
	Orange  Color = "\033[38;5;208m" // Orange color (using extended color)
)

// NewColored wraps an existing Logger and applies ANSI color to its messages.
func NewColored(base Logger, color Color) Logger {
	return &colored{
		base:  base,
		color: color,
	}
}

func (c *colored) Info(msg string, args ...any) {
	c.base.Info(c.colorize(msg), args...)
}

func (c *colored) Debug(msg string, args ...any) {
	c.base.Debug(c.colorize(msg), args...)
}

func (c *colored) Warn(msg string, args ...any) {
	c.base.Warn(c.colorize(msg), args...)
}

func (c *colored) Error(msg string, args ...any) {
	c.base.Error(c.colorize(msg), args...)
}

func (c *colored) colorize(msg string) string {
	return fmt.Sprintf("%s%s\033[0m", c.color, msg)
}
