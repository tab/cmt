package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParsePrefix(t *testing.T) {
	c := &commitCmd{}

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Success without prefix argument",
			args:     []string{},
			expected: "",
		},
		{
			name:     "Success with long flag value",
			args:     []string{"--prefix", "feat:"},
			expected: "feat:",
		},
		{
			name:     "Success with short flag value",
			args:     []string{"-p", "fix:"},
			expected: "fix:",
		},
		{
			name:     "Success with literal command value",
			args:     []string{"prefix", "chore:"},
			expected: "chore:",
		},
		{
			name:     "Success with inline long flag",
			args:     []string{"--prefix=docs:"},
			expected: "docs:",
		},
		{
			name:     "Success with inline short flag",
			args:     []string{"-p=refactor:"},
			expected: "refactor:",
		},
		{
			name:     "Success with uppercase flag",
			args:     []string{"--PREFIX=test:"},
			expected: "test:",
		},
		{
			name:     "Success with whitespace value",
			args:     []string{"--prefix", "  build:  "},
			expected: "build:",
		},
		{
			name:     "Success with inline whitespace value",
			args:     []string{"--prefix=  ci:  "},
			expected: "ci:",
		},
		{
			name:     "Success with prefix amid args",
			args:     []string{"some", "args", "--prefix", "perf:", "more", "args"},
			expected: "perf:",
		},
		{
			name:     "Success with inline prefix amid args",
			args:     []string{"some", "args", "--prefix=style:", "more", "args"},
			expected: "style:",
		},
		{
			name:     "Failure when value missing",
			args:     []string{"--prefix"},
			expected: "",
		},
		{
			name:     "Failure with following flag",
			args:     []string{"--prefix", "--some-flag"},
			expected: "",
		},
		{
			name:     "Success with mixed case flag",
			args:     []string{"--PrEfIx", "revert:"},
			expected: "revert:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.parsePrefix(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}
