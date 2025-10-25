package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cmt/internal/app/errors"
)

func Test_Format(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNoGitChanges",
			err:      errors.ErrNoGitChanges,
			expected: "⚠️ no changes to commit",
		},
		{
			name:     "ErrNoResponse",
			err:      errors.ErrNoResponse,
			expected: "⚠️ no response from GPT",
		},
		{
			name:     "ErrFailedToParseJSON",
			err:      errors.ErrFailedToParseJSON,
			expected: "⚠️ failed to parse JSON response",
		},
		{
			name:     "ErrCommitMessageEmpty",
			err:      errors.ErrCommitMessageEmpty,
			expected: "⚠️ commit message is empty, commit aborted",
		},
		{
			name:     "Generic error",
			err:      errors.New("some error"),
			expected: "❌ some error",
		},
		{
			name:     "ErrAPITokenNotSet",
			err:      errors.ErrAPITokenNotSet,
			expected: "❌ API token not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Format(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
