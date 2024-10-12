package config

import (
	"os"
	"testing"

	"cmt/internal/errors"

	"github.com/stretchr/testify/assert"
)

func Test_GetAPIToken(t *testing.T) {
	type result struct {
		error error
		token string
	}
	tests := []struct {
		name     string
		env      string
		expected result
	}{
		{
			name: "API token set",
			env:  "test-api-token",
			expected: result{
				error: nil,
				token: "test-api-token",
			},
		},
		{
			name: "API token not set",
			env:  "",
			expected: result{
				error: errors.ErrAPITokenNotSet,
				token: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("OPENAI_API_KEY", tt.env)

			token, err := GetAPIToken()

			if tt.expected.error != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expected.error.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.token, token)
			}

			os.Unsetenv("OPENAI_API_KEY")
		})
	}
}
