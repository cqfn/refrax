package log

import "os"

type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

var single Logger = NewZerolog(os.Stdout, "info")

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func Set(logger Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}
	single = logger
}

func Default() Logger {
	if single == nil {
		panic("logger is not set")
	}
	return single
}
