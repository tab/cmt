package commit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/commands"
	"cmt/internal/git"
	"cmt/internal/gpt"
)

func Test_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitClient := git.NewMockGitClient(ctrl)
	mockGPTModelClient := gpt.NewMockGPTModelClient(ctrl)

	options := commands.GenerateOptions{
		Ctx:    context.Background(),
		Client: mockGitClient,
		Model:  mockGPTModelClient,
	}

	type result struct {
		output string
		err    bool
	}

	tests := []struct {
		name     string
		input    func() (string, error)
		before   func()
		expected result
	}{
		{
			name: "Success",
			input: func() (string, error) {
				return "y", nil
			},
			before: func() {
				mockGitClient.EXPECT().Diff(gomock.Any(), nil).Return("mock diff output", nil)
				mockGPTModelClient.EXPECT().FetchCommitMessage(gomock.Any(), "mock diff output").Return("feat(core): Description", nil)
				mockGitClient.EXPECT().Commit(gomock.Any(), "feat(core): Description").Return("[feature/core 29ca12d] feat(core): Description", nil)
			},
			expected: result{
				output: "🚀 Changes committed:\n[feature/core 29ca12d] feat(core): Description",
				err:    false,
			},
		},
		{
			name: "Commit aborted",
			input: func() (string, error) {
				return "n", nil
			},
			before: func() {
				mockGitClient.EXPECT().Diff(gomock.Any(), nil).Return("mock diff output", nil)
				mockGPTModelClient.EXPECT().FetchCommitMessage(gomock.Any(), "mock diff output").Return("feat(core): Description", nil)
			},
			expected: result{
				output: "❌ Commit aborted",
				err:    false,
			},
		},
		{
			name: "Error fetching diff",
			input: func() (string, error) {
				return "", nil
			},
			before: func() {
				mockGitClient.EXPECT().Diff(gomock.Any(), nil).Return("", fmt.Errorf("git diff error"))
			},
			expected: result{
				output: "❌ Error getting git diff: git diff error",
				err:    true,
			},
		},
		{
			name: "Error committing changes",
			input: func() (string, error) {
				return "y", nil
			},
			before: func() {
				mockGitClient.EXPECT().Diff(gomock.Any(), nil).Return("mock diff output", nil)
				mockGPTModelClient.EXPECT().FetchCommitMessage(gomock.Any(), "mock diff output").Return("feat(core): Description", nil)
				mockGitClient.EXPECT().Commit(gomock.Any(), "feat(core): Description").Return("", fmt.Errorf("git commit error"))
			},
			expected: result{
				output: "❌ Error committing changes: git commit error",
				err:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			r, w, _ := os.Pipe()
			defer r.Close()
			defer w.Close()
			origStdout := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = origStdout }()

			err := NewCommand(options, func() (string, error) { return tt.input() }).Generate()

			w.Close()
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			if tt.expected.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Contains(t, output, tt.expected.output)
		})
	}
}
