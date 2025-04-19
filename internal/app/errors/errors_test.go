package errors_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"cmt/internal/app/errors"
)

func Test_HandleDiffError(t *testing.T) {
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

func Test_HandleGitLogError(t *testing.T) {
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

func Test_HandleModelError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNoResponse",
			err:      errors.ErrNoResponse,
			expected: "⚠️ No response from GPT\n",
		},
		{
			name:     "ErrFailedToParseJSON",
			err:      errors.ErrFailedToParseJSON,
			expected: "⚠️ Failed to parse JSON response\n",
		},
		{
			name:     "Other error",
			err:      errors.New("some other error"),
			expected: "❌ Error getting model response: some other error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			errors.HandleModelError(tt.err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}

func Test_HandleInputError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrWrongInput",
			err:      errors.ErrWrongInput,
			expected: "⚠️ Invalid input, please enter 'y', 'e' or 'n'\n",
		},
		{
			name:     "Other error",
			err:      errors.New("some other error"),
			expected: "❌ Error reading user input: some other error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			errors.HandleInputError(tt.err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}

func Test_HandleCommitError(t *testing.T) {
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

func Test_HandleEditError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "Error",
			err:      errors.New("some error"),
			expected: "❌ Error editing commit message: some error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			errors.HandleEditError(tt.err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}
