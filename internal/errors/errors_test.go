package errors_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"cmt/internal/errors"
)

func TestHandleDiffError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNoGitChanges",
			err:      errors.ErrNoGitChanges,
			expected: "⚠️ No changes to commit\n",
		},
		{
			name:     "Other error",
			err:      errors.New("some other error"),
			expected: "❌ Error getting git diff: some other error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			errors.HandleDiffError(tt.err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestHandleGitLogError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNoGitChanges",
			err:      errors.ErrNoGitChanges,
			expected: "⚠️ No changes found in the git log\n",
		},
		{
			name:     "Other error",
			err:      errors.New("some other error"),
			expected: "❌ Error getting git log: some other error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			errors.HandleGitLogError(tt.err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestHandleCommitError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrCommitMessageEmpty",
			err:      errors.ErrCommitMessageEmpty,
			expected: "⚠️ Commit message is empty, commit aborted\n",
		},
		{
			name:     "Other error",
			err:      errors.New("some other error"),
			expected: "❌ Error committing changes: some other error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			errors.HandleCommitError(tt.err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}
