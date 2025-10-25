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
}

func Test_Load(t *testing.T) {
	oldToken := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", oldToken)

	os.Setenv("OPENAI_API_KEY", "test-token")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Model.Name)
	assert.Greater(t, cfg.Model.MaxTokens, 0)
	assert.GreaterOrEqual(t, cfg.API.RetryCount, 0)
	assert.NotEmpty(t, cfg.Logging.Level)
}

func Test_Load_WithInvalidConfig(t *testing.T) {
	oldToken := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", oldToken)

	os.Setenv("OPENAI_API_KEY", "test-token")

	configContent := `model:
  temperature: 3.0`

	err := os.WriteFile("cmt.yaml", []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.Remove("cmt.yaml")

	cfg, err := Load()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func Test_Load_WithBadYAML(t *testing.T) {
	oldToken := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", oldToken)

	os.Setenv("OPENAI_API_KEY", "test-token")

	configContent := `model:
  name: [invalid yaml structure`

	err := os.WriteFile("cmt.yaml", []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.Remove("cmt.yaml")

	cfg, err := Load()

	assert.Error(t, err)
	assert.Nil(t, cfg)
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

func Test_Validate(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config - all values correct",
			setupConfig: DefaultConfig,
			expectError: false,
		},
		{
			name: "valid config - temperature at lower bound",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = 0
				return cfg
			},
			expectError: false,
		},
		{
			name: "valid config - temperature at upper bound",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = 2
				return cfg
			},
			expectError: false,
		},
		{
			name: "valid config - retry count zero",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.RetryCount = 0
				return cfg
			},
			expectError: false,
		},
		{
			name: "invalid temperature - negative",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = -0.1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid temperature",
		},
		{
			name: "invalid temperature - too high",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = 2.1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid temperature",
		},
		{
			name: "invalid max_tokens - zero",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.MaxTokens = 0
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid max_tokens",
		},
		{
			name: "invalid max_tokens - negative",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.MaxTokens = -100
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid max_tokens",
		},
		{
			name: "invalid timeout - zero",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.Timeout = 0
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid timeout",
		},
		{
			name: "invalid timeout - negative",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.Timeout = -1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid timeout",
		},
		{
			name: "invalid retry_count - negative",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.RetryCount = -1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid retry_count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setupConfig()
			err := cfg.validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
