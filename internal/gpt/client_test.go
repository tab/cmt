package gpt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	tests := []struct {
		name   string
		before func()
		token  string
		error  bool
	}{
		{
			name: "Success",
			before: func() {
				os.Setenv("OPENAI_API_KEY", "secret-token")
			},
			error: false,
		},
		{
			name: "Error",
			before: func() {
				os.Setenv("OPENAI_API_KEY", "")
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			result, err := NewClient()

			if tt.error {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			os.Unsetenv("OPENAI_API_KEY")
		})
	}
}
