package log

import (
	"io"

	"github.com/rs/zerolog"
)

type zero struct {
	logger zerolog.Logger
}

// NewZerolog creates a new Logger with the specified log level and writer.
// Details: https://github.com/rs/zerolog
func NewZerolog(writer io.Writer, level string) Logger {
	zlevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zlevel = zerolog.InfoLevel
	}
	return &zero{
		logger: zerolog.New(zerolog.ConsoleWriter{Out: writer}).Level(zlevel).With().Timestamp().Logger(),
	}
}

func (z *zero) Info(msg string, args ...any) {
	z.logger.Info().Msgf(msg, args...)
}

func (z *zero) Debug(msg string, args ...any) {
	z.logger.Debug().Msgf(msg, args...)
}

func (z *zero) Warn(msg string, args ...any) {
	z.logger.Warn().Msgf(msg, args...)
}

func (z *zero) Error(msg string, args ...any) {
	z.logger.Error().Msgf(msg, args...)
}
