package git

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GitDiff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	g := &Git{
		Executor: mockExecutor,
	}

	type result struct {
		output string
		err    error
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--minimal", "--ignore-all-space", "--ignore-blank-lines", "--staged").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "mock diff output")
							cmd.Stdout = bytes.NewBufferString("mock diff output")
							return cmd
						})
			},
			expected: result{
				output: "mock diff output",
				err:    nil,
			},
		},
		{
			name: "Failure",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--minimal", "--ignore-all-space", "--ignore-blank-lines", "--staged").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("false")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.New("failed to load git diff"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			result, err := g.Diff(context.Background(), nil)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, result)
			}
		})
	}
}

func TestGit_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	g := &Git{
		Executor: mockExecutor,
	}

	type result struct {
		output string
		err    error
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "log", "--format='%h %s %b'", "v1.0.0..v1.2.0").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "mock log output")
							cmd.Stdout = bytes.NewBufferString("mock log output")
							return cmd
						})
			},
			expected: result{
				output: "mock log output",
				err:    nil,
			},
		},
		{
			name: "Failure",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "log", "--format='%h %s %b'", "v1.0.0..v1.2.0").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("false")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.New("failed to load git log"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			result, err := g.Log(context.Background(), []string{"v1.0.0..v1.2.0"})

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, result)
			}
		})
	}

}

func Test_GitCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	g := &Git{
		Executor: mockExecutor,
	}

	type result struct {
		output string
		err    error
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "commit", "-m", "test commit").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "commit successful")
							cmd.Stdout = bytes.NewBufferString("mock commit output")
							return cmd
						})
			},
			expected: result{
				output: "commit successful\n",
				err:    nil,
			},
		},
		{
			name: "Failure",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "commit", "-m", "test commit").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("false")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.New("failed to commit changes"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			result, err := g.Commit(context.Background(), "test commit")

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, result)
			}
		})
	}
}

type fakeExecutor struct{}

func (f *fakeExecutor) Run(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, arg...)
}

func Test_GitEdit(t *testing.T) {
	tempDir := t.TempDir()
	scriptPath := filepath.Join(tempDir, "fake_editor.sh")
	scriptContent := `#!/bin/sh
echo "edited message" > "$1"
`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	assert.NoError(t, err)

	err = os.Setenv("EDITOR", scriptPath)
	assert.NoError(t, err)
	defer os.Unsetenv("EDITOR")

	g := &Git{Executor: &fakeExecutor{}}

	ctx := context.Background()
	originalMessage := "original message"
	edited, err := g.Edit(ctx, originalMessage)
	assert.NoError(t, err)
	assert.Equal(t, "edited message", edited)
}
