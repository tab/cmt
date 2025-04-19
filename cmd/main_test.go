package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cmt/internal/config"
)

func Test_Main(t *testing.T) {
	cleanup := generateConfigFile(t)
	defer cleanup()

	tests := []struct {
		name  string
		env   map[string]string
		error bool
	}{
		{
			name: "Success",
			env: map[string]string{
				"OPENAI_API_KEY": "test-api-key-1234",
			},
			error: false,
		},
		{
			name: "Missing API Key",
			env: map[string]string{
				"OPENAI_API_KEY": "",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEnv := make(map[string]string)
			for k := range tt.env {
				originalEnv[k] = os.Getenv(k)
			}

			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			defer func() {
				for k, v := range originalEnv {
					os.Setenv(k, v)
				}
			}()

			cfg, err := config.Load()
			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

func generateConfigFile(t *testing.T) func() {
	t.Helper()

	filename := "cmt.yaml"
	content := "editor: vim\n"

	originalContent := []byte{}
	fileExists := false
	if _, err := os.Stat(filename); err == nil {
		fileExists = true
		originalContent, err = os.ReadFile(filename)
		require.NoError(t, err)
	}

	err := os.WriteFile(filename, []byte(content), 0644)
	require.NoError(t, err)

	_, err = os.ReadFile(filename)
	require.NoError(t, err)

	return func() {
		if fileExists {
			err = os.WriteFile(filename, originalContent, 0644)
			if err != nil {
				t.Logf("Warning: Failed to restore original config file: %v", err)
			}
		} else {
			err = os.Remove(filename)
			if err != nil {
				t.Logf("Warning: Failed to remove test config file: %v", err)
			}
		}
	}
}
