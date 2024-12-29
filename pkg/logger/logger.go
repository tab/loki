package logger

import (
	"io"
	"os"

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
		switch envValue {
		case "debug":
			return zerolog.DebugLevel
		case "info":
			return zerolog.InfoLevel
		case "warn":
			return zerolog.WarnLevel
		case "error":
			return zerolog.ErrorLevel
		case "fatal":
			return zerolog.FatalLevel
		case "panic":
			return zerolog.PanicLevel
		case "trace":
			return zerolog.TraceLevel
		default:
			return zerolog.InfoLevel
		}
	}

	return zerolog.InfoLevel
}
