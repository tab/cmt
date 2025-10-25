package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertCommitLogForGPT(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty log",
			input:    "",
			expected: "",
		},
		{
			name:     "single commit",
			input:    "abc123|feat: add feature|john|2 days ago",
			expected: "abc123 feat: add feature",
		},
		{
			name:     "multiple commits",
			input:    "abc123|feat: add feature|john|2 days ago\ndef456|fix: bug fix|jane|3 days ago",
			expected: "abc123 feat: add feature\ndef456 fix: bug fix",
		},
		{
			name:     "malformed line",
			input:    "abc123",
			expected: "",
		},
		{
			name:     "mixed valid and invalid",
			input:    "abc123|feat: add feature|john|2 days ago\ninvalid\ndef456|fix: bug fix|jane|3 days ago",
			expected: "abc123 feat: add feature\ndef456 fix: bug fix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCommitLogForGPT(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
