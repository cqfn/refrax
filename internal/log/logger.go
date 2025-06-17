package log

import "os"

type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

var single Logger = NewZerolog(os.Stdout, "debug")

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

func Set(logger Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}
	single = logger
}

func Get() Logger {
	if single == nil {
		panic("logger is not set")
	}
	return single
}
