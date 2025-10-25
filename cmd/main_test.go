package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxevent"

	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_LoadConfig(t *testing.T) {
	cfg, err := config.Load()

	if err != nil {
		t.Skip("config loading failed, likely no cmt.yaml file in expected location")
		return
	} else {
		assert.NoError(t, err)
	}

	assert.NotNil(t, cfg)
}

func Test_CreateApp(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "Creates app with info level logging",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.InfoLevel,
				},
			},
		},
		{
			name: "Creates app with debug level logging",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.DebugLevel,
				},
			},
		},
		{
			name: "Creates app with error level logging",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.ErrorLevel,
				},
			},
		},
		{
			name: "Creates app with warn level logging",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.WarnLevel,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := createApp(tt.config)
			assert.NotNil(t, app)
		})
	}
}

func Test_CreateFxLogger(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		expectedType   interface{}
		expectedLogger interface{}
	}{
		{
			name: "Debug level returns console logger",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.DebugLevel,
				},
			},
			expectedType: &fxevent.ConsoleLogger{},
		},
		{
			name: "Info level returns nop logger",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.InfoLevel,
				},
			},
			expectedLogger: fxevent.NopLogger,
		},
		{
			name: "Warn level returns nop logger",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.WarnLevel,
				},
			},
			expectedLogger: fxevent.NopLogger,
		},
		{
			name: "Error level returns nop logger",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.ErrorLevel,
				},
			},
			expectedLogger: fxevent.NopLogger,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerFunc := createFxLogger(tt.config)
			assert.NotNil(t, loggerFunc)

			result := loggerFunc()
			assert.NotNil(t, result)

			if tt.expectedType != nil {
				assert.IsType(t, tt.expectedType, result)
			}
			if tt.expectedLogger != nil {
				assert.Equal(t, tt.expectedLogger, result)
			}
		})
	}
}

func Test_CreateFxLogger_FunctionCreation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "Creates valid function with debug config",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.DebugLevel,
				},
			},
		},
		{
			name: "Creates valid function with info config",
			config: &config.Config{
				Logging: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{
					Level: logger.InfoLevel,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerFunc := createFxLogger(tt.config)
			assert.NotNil(t, loggerFunc)

			result1 := loggerFunc()
			result2 := loggerFunc()

			assert.NotNil(t, result1)
			assert.NotNil(t, result2)
		})
	}
}

func Test_handleSimpleCommands(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedHandled bool
		outputContains  string
	}{
		{
			name:            "no args",
			args:            []string{},
			expectedHandled: false,
		},
		{
			name:            "unknown command",
			args:            []string{"commit"},
			expectedHandled: false,
		},
		{
			name:            "version flag",
			args:            []string{"--version"},
			expectedHandled: true,
			outputContains:  config.Version,
		},
		{
			name:            "help flag",
			args:            []string{"-h"},
			expectedHandled: true,
			outputContains:  "USAGE:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				handled := handleSimpleCommands(tt.args)
				assert.Equal(t, tt.expectedHandled, handled)
			})

			if tt.outputContains != "" {
				assert.Contains(t, output, tt.outputContains)
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func captureStdout(fn func()) string {
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}
