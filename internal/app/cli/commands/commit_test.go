package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePrefix(t *testing.T) {
	c := &commitCmd{}

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "no prefix argument",
			args:     []string{},
			expected: "",
		},
		{
			name:     "--prefix with space-separated value",
			args:     []string{"--prefix", "feat:"},
			expected: "feat:",
		},
		{
			name:     "-p with space-separated value",
			args:     []string{"-p", "fix:"},
			expected: "fix:",
		},
		{
			name:     "prefix with space-separated value",
			args:     []string{"prefix", "chore:"},
			expected: "chore:",
		},
		{
			name:     "--prefix=value inline format",
			args:     []string{"--prefix=docs:"},
			expected: "docs:",
		},
		{
			name:     "-p=value inline format",
			args:     []string{"-p=refactor:"},
			expected: "refactor:",
		},
		{
			name:     "--PREFIX=value uppercase",
			args:     []string{"--PREFIX=test:"},
			expected: "test:",
		},
		{
			name:     "prefix with extra whitespace",
			args:     []string{"--prefix", "  build:  "},
			expected: "build:",
		},
		{
			name:     "prefix=value with extra whitespace",
			args:     []string{"--prefix=  ci:  "},
			expected: "ci:",
		},
		{
			name:     "prefix among other arguments",
			args:     []string{"some", "args", "--prefix", "perf:", "more", "args"},
			expected: "perf:",
		},
		{
			name:     "inline prefix among other arguments",
			args:     []string{"some", "args", "--prefix=style:", "more", "args"},
			expected: "style:",
		},
		{
			name:     "prefix without value",
			args:     []string{"--prefix"},
			expected: "",
		},
		{
			name:     "prefix with flag as next value (should not use it)",
			args:     []string{"--prefix", "--some-flag"},
			expected: "",
		},
		{
			name:     "mixed case prefix",
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
