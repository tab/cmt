package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, DefaultModelName, cfg.Model.Name)
	assert.Equal(t, DefaultMaxTokens, cfg.Model.MaxTokens)
	assert.Equal(t, DefaultTemperature, cfg.Model.Temperature)
	assert.Equal(t, DefaultRetryCount, cfg.API.RetryCount)
	assert.Equal(t, DefaultTimeout, cfg.API.Timeout)
	assert.Equal(t, DefaultLogLevel, cfg.Logging.Level)
	assert.Equal(t, DefaultLogFormat, cfg.Logging.Format)
	assert.Equal(t, DefaultEditor, cfg.Editor)
}

func Test_Load(t *testing.T) {
	tests := []struct {
		name          string
		configFile    string
		configContent string
		token         string
		error         bool
		expected      *Config
	}{
		{
			name:     "Default config",
			token:    "test-token",
			error:    false,
			expected: DefaultConfig(),
		},
		{
			name:  "Missing api token",
			token: "",
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldToken := os.Getenv("OPENAI_API_KEY")
			defer os.Setenv("OPENAI_API_KEY", oldToken)

			os.Setenv("OPENAI_API_KEY", tt.token)

			cfg, err := Load()

			if tt.error {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, tt.expected.Model.Name, cfg.Model.Name)
				assert.Equal(t, tt.expected.Model.MaxTokens, cfg.Model.MaxTokens)
				assert.Equal(t, tt.expected.API.RetryCount, cfg.API.RetryCount)
				assert.Equal(t, tt.expected.Logging.Level, cfg.Logging.Level)
			}
		})
	}
}

func Test_GetAPIToken(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		error bool
	}{
		{
			name:  "Success",
			env:   "test-token-12345",
			error: false,
		},
		{
			name:  "Error",
			env:   "",
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldEnv := os.Getenv("OPENAI_API_KEY")
			defer os.Setenv("OPENAI_API_KEY", oldEnv)

			os.Setenv("OPENAI_API_KEY", tt.env)

			token, err := GetAPIToken()
			if tt.error {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.env, token)
			}
		})
	}
}
