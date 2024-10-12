package git

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
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
				mockExecutor.EXPECT().Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").Return(nil).DoAndReturn(func(ctx context.Context, name string, args ...string) *exec.Cmd {
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
				mockExecutor.EXPECT().Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").Return(nil).DoAndReturn(func(ctx context.Context, name string, args ...string) *exec.Cmd {
					cmd := exec.Command("false")
					return cmd
				})
			},
			expected: result{
				output: "",
				err:    errors.New("git diff error: exit status 1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			result, err := g.Diff(context.Background())

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
				mockExecutor.EXPECT().Run(gomock.Any(), "git", "commit", "-m", "test commit").Return(nil).DoAndReturn(func(ctx context.Context, name string, args ...string) *exec.Cmd {
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
				mockExecutor.EXPECT().Run(gomock.Any(), "git", "commit", "-m", "test commit").Return(nil).DoAndReturn(func(ctx context.Context, name string, args ...string) *exec.Cmd {
					cmd := exec.Command("false")
					return cmd
				})
			},
			expected: result{
				output: "",
				err:    errors.New("git commit error: exit status 1"),
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
