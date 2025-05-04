package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"loki/internal/config"
)

func Test_NewLogger(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		expected zerolog.Level
	}{
		{
			name: "Debug level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: DebugLevel,
			},
			expected: zerolog.DebugLevel,
		},
		{
			name: "Info level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: InfoLevel,
			},
			expected: zerolog.InfoLevel,
		},
		{
			name: "Warn level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: WarnLevel,
			},
			expected: zerolog.WarnLevel,
		},
		{
			name: "Error level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: ErrorLevel,
			},
			expected: zerolog.ErrorLevel,
		},
		{
			name: "Fatal level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: FatalLevel,
			},
			expected: zerolog.FatalLevel,
		},
		{
			name: "Panic level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: PanicLevel,
			},
			expected: zerolog.PanicLevel,
		},
		{
			name: "Trace level",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: TraceLevel,
			},
			expected: zerolog.TraceLevel,
		},
		{
			name: "Default level when empty",
			cfg: &config.Config{
				AppName:  "test-app",
				AppEnv:   "test",
				LogLevel: "",
			},
			expected: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.cfg)

			assert.Equal(t, tt.expected, logger.log.GetLevel())
			assert.NotNil(t, logger)
		})
	}
}

func Test_Logger_WithComponent(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "debug",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	componentLogger := logger.WithComponent("test-component")
	componentLogger.Info().Msg("test message")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "test-component", logData["component"])
	assert.Equal(t, "test message", logData["message"])
	assert.Equal(t, "test-app", logData["service"])
	assert.Equal(t, "test", logData["environment"])
}

func Test_Logger_WithRequestId(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "debug",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	requestLogger := logger.WithRequestId("req-123")
	requestLogger.Info().Msg("handling request")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "req-123", logData["request_id"])
}

func Test_Logger_WithTraceId(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "debug",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	traceLogger := logger.WithTraceId("trace-123")
	traceLogger.Info().Msg("traced operation")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "trace-123", logData["trace_id"])
}

func Test_Logger_Debug(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "debug",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	logger.Debug().Msg("debug message")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "debug message", logData["message"])
}

func Test_Logger_Info(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "info",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	logger.Info().Msg("info message")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "info message", logData["message"])
}

func Test_Logger_Warn(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "warn",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	logger.Warn().Msg("warn message")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "warn message", logData["message"])
}

func Test_Logger_Error(t *testing.T) {
	var buf bytes.Buffer

	cfg := &config.Config{
		AppName:  "test-app",
		AppEnv:   "test",
		LogLevel: "error",
	}

	logger := NewLogger(cfg)
	logger.log = logger.log.Output(&buf)

	logger.Error().Msg("error message")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)

	assert.NoError(t, err)
	assert.Equal(t, "error message", logData["message"])
}
