package log

import (
	"io"

	"github.com/rs/zerolog"
)

type Zerolog struct {
	logger zerolog.Logger
}

func NewZerolog(writer io.Writer, level string) Logger {
	zlevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zlevel = zerolog.InfoLevel
	}
	return &Zerolog{
		logger: zerolog.New(zerolog.ConsoleWriter{Out: writer}).Level(zlevel).With().Timestamp().Logger(),
	}
}

func (z *Zerolog) Info(msg string, args ...any) {
	z.logger.Info().Msgf(msg, args...)
}

func (z *Zerolog) Debug(msg string, args ...any) {
	z.logger.Debug().Msgf(msg, args...)
}

func (z *Zerolog) Warn(msg string, args ...any) {
	z.logger.Warn().Msgf(msg, args...)
}

func (z *Zerolog) Error(msg string, args ...any) {
	z.logger.Error().Msgf(msg, args...)
}
