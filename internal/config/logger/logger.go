package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"loki/internal/config"
)

const (
	Component = "component"
	RequestId = "request_id"
	TraceId   = "trace_id"

	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
	TraceLevel = "trace"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger(cfg *config.Config) *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"

	var output io.Writer = os.Stdout

	hostname, _ := os.Hostname()

	log := zerolog.New(output).
		Level(getLogLevel(cfg.LogLevel)).
		With().
		Timestamp().
		Str("service", cfg.AppName).
		Str("environment", cfg.AppEnv).
		Str("host", hostname).
		Logger()

	return &Logger{log: log}
}

func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		log: l.log.With().Str(Component, component).Logger(),
	}
}

func (l *Logger) WithRequestId(requestId string) *Logger {
	return &Logger{
		log: l.log.With().Str(RequestId, requestId).Logger(),
	}
}

func (l *Logger) WithTraceId(traceId string) *Logger {
	return &Logger{
		log: l.log.With().Str(TraceId, traceId).Logger(),
	}
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

func getLogLevel(level string) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	case TraceLevel:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}
