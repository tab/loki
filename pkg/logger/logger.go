package logger

import (
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	var output io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: zerolog.TimeFieldFormat,
	}

	log := zerolog.New(output).
		Level(getLogLevel()).
		With().
		Timestamp().
		Logger()

	return &Logger{log: log}
}

func (l *Logger) Debug() *zerolog.Event {
	return l.log.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.log.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.log.Error()
}

func getLogLevel() zerolog.Level {
	if envValue, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := strconv.Atoi(envValue); err == nil && level >= 0 && level <= 5 {
			return zerolog.Level(level)
		}
	}

	return zerolog.InfoLevel
}
