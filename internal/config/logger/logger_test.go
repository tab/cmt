package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"cmt/internal/config"
)

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_NewLogger(t *testing.T) {
	type result struct {
		level zerolog.Level
	}

	tests := []struct {
		name     string
		cfg      *config.Config
		expected result
	}{
		{
			name: "Success with default level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				return cfg
			}(),
			expected: result{
				level: zerolog.InfoLevel,
			},
		},
		{
			name: "Success with empty level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = ""
				return cfg
			}(),
			expected: result{
				level: zerolog.InfoLevel,
			},
		},
		{
			name: "Success with debug level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = DebugLevel
				return cfg
			}(),
			expected: result{
				level: zerolog.DebugLevel,
			},
		},
		{
			name: "Success with warn level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = WarnLevel
				return cfg
			}(),
			expected: result{
				level: zerolog.WarnLevel,
			},
		},

		{
			name: "Success with error level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = ErrorLevel
				return cfg
			}(),
			expected: result{
				level: zerolog.ErrorLevel,
			},
		},
		{
			name: "Success with fatal level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = FatalLevel
				return cfg
			}(),
			expected: result{
				level: zerolog.FatalLevel,
			},
		},
		{
			name: "Success with panic level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = PanicLevel
				return cfg
			}(),
			expected: result{
				level: zerolog.PanicLevel,
			},
		},
		{
			name: "Success with trace level",
			cfg: func() *config.Config {
				cfg := config.DefaultConfig()
				cfg.Logging.Level = TraceLevel
				return cfg
			}(),
			expected: result{
				level: zerolog.TraceLevel,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.cfg)
			assert.NotNil(t, logger)

			appLogger, ok := logger.(*AppLogger)
			assert.True(t, ok)

			assert.Equal(t, tt.expected.level, appLogger.log.GetLevel())
		})
	}
}

func Test_Logger_Debug(t *testing.T) {
	cfg := &config.Config{
		Logging: struct {
			Level string `yaml:"level"`
		}{
			Level: DebugLevel,
		},
	}

	logger := NewLogger(cfg)
	logger.Debug().Msg("debug message")

	assert.NotNil(t, logger)
}

func Test_Logger_Info(t *testing.T) {
	cfg := &config.Config{
		Logging: struct {
			Level string `yaml:"level"`
		}{
			Level: InfoLevel,
		},
	}

	logger := NewLogger(cfg)
	logger.Info().Msg("info message")

	assert.NotNil(t, logger)
}

func Test_Logger_Warn(t *testing.T) {
	cfg := &config.Config{
		Logging: struct {
			Level string `yaml:"level"`
		}{
			Level: WarnLevel,
		},
	}

	logger := NewLogger(cfg)
	logger.Warn().Msg("warn message")

	assert.NotNil(t, logger)
}

func Test_Logger_Error(t *testing.T) {
	cfg := &config.Config{
		Logging: struct {
			Level string `yaml:"level"`
		}{
			Level: ErrorLevel,
		},
	}

	logger := NewLogger(cfg)
	logger.Error().Msg("error message")

	assert.NotNil(t, logger)
}

func Test_GetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected zerolog.Level
	}{
		{
			name:     "Success with debug level",
			level:    DebugLevel,
			expected: zerolog.DebugLevel,
		},
		{
			name:     "Success with info level",
			level:    InfoLevel,
			expected: zerolog.InfoLevel,
		},
		{
			name:     "Success with warn level",
			level:    WarnLevel,
			expected: zerolog.WarnLevel,
		},
		{
			name:     "Success with error level",
			level:    ErrorLevel,
			expected: zerolog.ErrorLevel,
		},
		{
			name:     "Success with fatal level",
			level:    FatalLevel,
			expected: zerolog.FatalLevel,
		},
		{
			name:     "Success with panic level",
			level:    PanicLevel,
			expected: zerolog.PanicLevel,
		},
		{
			name:     "Success with trace level",
			level:    TraceLevel,
			expected: zerolog.TraceLevel,
		},
		{
			name:     "Success with unknown level",
			level:    "unknown",
			expected: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getLogLevel(tt.level)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func Test_NewLogger_WithBuffer(t *testing.T) {
	cfg := &config.Config{
		Logging: struct {
			Level string `yaml:"level"`
		}{
			Level: InfoLevel,
		},
	}

	logger := NewLogger(cfg)

	assert.NotNil(t, logger)
	assert.NotNil(t, logger.GetBuffer())

	logger.Info().Int("diff_size", 8740).Str("version", "0.7.0").Msg("Generating commit message")

	buffer := logger.GetBuffer()
	entries := buffer.Entries()
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "info", entries[0].Level)
	assert.Contains(t, entries[0].Message, "Generating commit message")
	assert.Contains(t, entries[0].Message, "diff_size=8740")
	assert.Contains(t, entries[0].Message, "version=0.7.0")
}

func Test_LogBuffer(t *testing.T) {
	buffer := NewLogBuffer(3)

	buffer.Add("info", "message 1")
	buffer.Add("debug", "message 2")
	buffer.Add("error", "message 3")

	entries := buffer.Entries()
	assert.Equal(t, 3, len(entries))
	assert.Equal(t, "info", entries[0].Level)
	assert.Equal(t, "message 1", entries[0].Message)

	buffer.Add("warn", "message 4")

	entries = buffer.Entries()
	assert.Equal(t, 3, len(entries))
	assert.Equal(t, "debug", entries[0].Level)
	assert.Equal(t, "message 2", entries[0].Message)
	assert.Equal(t, "warn", entries[2].Level)
	assert.Equal(t, "message 4", entries[2].Message)

	buffer.Clear()
	entries = buffer.Entries()
	assert.Equal(t, 0, len(entries))
}
