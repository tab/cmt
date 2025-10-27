package errors_test

import (
	"fmt"
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
			name:     "Success with no git changes",
			err:      errors.ErrNoGitChanges,
			expected: "⚠️ no changes to commit",
		},
		{
			name:     "Success with no response",
			err:      errors.ErrNoResponse,
			expected: "⚠️ no response from GPT",
		},
		{
			name:     "Success with failed j s o n parse",
			err:      errors.ErrFailedToParseJSON,
			expected: "⚠️ failed to parse JSON response",
		},
		{
			name:     "Success with empty commit message",
			err:      errors.ErrCommitMessageEmpty,
			expected: "⚠️ commit message is empty, commit aborted",
		},
		{
			name:     "Success with unknown command",
			err:      errors.ErrUnknownCommand,
			expected: "⚠️ unknown command. Use 'cmt --help' for usage",
		},
		{
			name:     "Success with generic error",
			err:      errors.New("some error"),
			expected: "❌ some error",
		},
		{
			name:     "Success with missing a p i token",
			err:      errors.ErrAPITokenNotSet,
			expected: "❌ API token not set",
		},
		{
			name:     "Success with wrapped no git changes",
			err:      fmt.Errorf("context: %w", errors.ErrNoGitChanges),
			expected: "⚠️ no changes to commit",
		},
		{
			name:     "Success with wrapped generic error",
			err:      fmt.Errorf("wrapped: %w", errors.New("original error")),
			expected: "❌ wrapped: original error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Format(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
