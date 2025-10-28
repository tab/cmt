package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// writeTempConfig creates a temporary directory with a cmt.yaml file
// containing the specified contents and returns the directory path.
func writeTempConfig(t *testing.T, contents string) string {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cmt.yaml")
	err := os.WriteFile(configPath, []byte(contents), 0600)
	if err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return tmpDir
}

func Test_DefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, DefaultModelName, cfg.Model.Name)
	assert.Equal(t, DefaultMaxTokens, cfg.Model.MaxTokens)
	assert.Equal(t, DefaultTemperature, cfg.Model.Temperature)
	assert.Equal(t, DefaultRetryCount, cfg.API.RetryCount)
	assert.Equal(t, DefaultTimeout, cfg.API.Timeout)
	assert.Equal(t, DefaultLogLevel, cfg.Logging.Level)
}

func Test_Load(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-token")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Model.Name)
	assert.Greater(t, cfg.Model.MaxTokens, 0)
	assert.GreaterOrEqual(t, cfg.API.RetryCount, 0)
	assert.NotEmpty(t, cfg.Logging.Level)
}

func Test_Load_WithInvalidConfig(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-token")

	configContent := `model:
  temperature: 3.0`

	tmpDir := writeTempConfig(t, configContent)
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	cfg, err := Load()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func Test_Load_WithBadYAML(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-token")

	configContent := `model:
  name: [invalid yaml structure`

	tmpDir := writeTempConfig(t, configContent)
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

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
			name:  "Failure with missing token",
			env:   "",
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", tt.env)

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
			name:        "Success",
			setupConfig: DefaultConfig,
			expectError: false,
		},
		{
			name: "Success with temperature at lower bound",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = 0
				return cfg
			},
			expectError: false,
		},
		{
			name: "Success with temperature at upper bound",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = 2
				return cfg
			},
			expectError: false,
		},
		{
			name: "Success with zero retry count",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.RetryCount = 0
				return cfg
			},
			expectError: false,
		},
		{
			name: "Failure with negative temperature",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = -0.1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid temperature",
		},
		{
			name: "Failure with high temperature",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.Temperature = 2.1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid temperature",
		},
		{
			name: "Failure with zero max tokens",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.MaxTokens = 0
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid max_tokens",
		},
		{
			name: "Failure with negative max tokens",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.Model.MaxTokens = -100
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid max_tokens",
		},
		{
			name: "Failure with zero timeout",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.Timeout = 0
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid timeout",
		},
		{
			name: "Failure with negative timeout",
			setupConfig: func() *Config {
				cfg := DefaultConfig()
				cfg.API.Timeout = -1
				return cfg
			},
			expectError: true,
			errorMsg:    "invalid timeout",
		},
		{
			name: "Failure with negative retry count",
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
