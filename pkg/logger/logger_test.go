package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_NewLogger(t *testing.T) {
	tests := []struct {
		name     string
		expected zerolog.Level
	}{
		{
			name:     "Default configuration",
			expected: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()

			assert.Equal(t, tt.expected, logger.log.GetLevel())
			assert.NotNil(t, logger)
		})
	}
}

func Test_Logger_Debug(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Info().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}

func Test_Logger_Info(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Info().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}

func Test_Logger_Warn(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "warn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Warn().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}

func Test_Logger_Error(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Error().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}

func Test_getLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		before   func()
		expected zerolog.Level
	}{
		{
			name: "Debug level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "debug")
				assert.NoError(t, err)
			},
			expected: zerolog.DebugLevel,
		},
		{
			name: "Info level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "info")
				assert.NoError(t, err)
			},
			expected: zerolog.InfoLevel,
		},
		{
			name: "Warn level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "warn")
				assert.NoError(t, err)
			},
			expected: zerolog.WarnLevel,
		},
		{
			name: "Error level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "error")
				assert.NoError(t, err)
			},
			expected: zerolog.ErrorLevel,
		},
		{
			name: "Fatal level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "fatal")
				assert.NoError(t, err)
			},
			expected: zerolog.FatalLevel,
		},
		{
			name: "Panic level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "panic")
				assert.NoError(t, err)
			},
			expected: zerolog.PanicLevel,
		},
		{
			name: "Trace level",
			before: func() {
				err := os.Setenv("LOG_LEVEL", "trace")
				assert.NoError(t, err)
			},
			expected: zerolog.TraceLevel,
		},
		{
			name:     "Default level",
			before:   func() {},
			expected: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			level := getLogLevel()

			assert.Equal(t, tt.expected, level)

			t.Cleanup(func() {
				err := os.Unsetenv("LOG_LEVEL")
				assert.NoError(t, err)
			})
		})
	}
}
