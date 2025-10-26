package git

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/errors"
	"cmt/internal/config/logger"
)

type stubLogger struct {
	log zerolog.Logger
}

func newStubLogger() stubLogger {
	return stubLogger{
		log: zerolog.Nop(),
	}
}

func (s stubLogger) Debug() *zerolog.Event        { return s.log.Debug() }
func (s stubLogger) Info() *zerolog.Event         { return s.log.Info() }
func (s stubLogger) Warn() *zerolog.Event         { return s.log.Warn() }
func (s stubLogger) Error() *zerolog.Event        { return s.log.Error() }
func (s stubLogger) GetBuffer() *logger.LogBuffer { return nil }

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_NewGitClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	log := newStubLogger()

	gitClient := NewGitClient(mockExecutor, log)
	assert.NotNil(t, gitClient)

	instance, ok := gitClient.(*client)
	assert.True(t, ok)
	assert.NotNil(t, instance.executor)
	assert.NotNil(t, instance.log)
}

func Test_Diff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockExecutor := NewMockExecutor(ctrl)
	log := newStubLogger()

	gitClient := &client{
		executor: mockExecutor,
		log:      log,
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
					Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").
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
					Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").
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
		{
			name: "No changes to commit",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "")
							cmd.Stdout = bytes.NewBufferString("")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.ErrNoGitChanges,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			output, err := gitClient.Diff(ctx)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, output)
			}
		})
	}
}

func Test_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	log := newStubLogger()

	gitClient := &client{
		executor: mockExecutor,
		log:      log,
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
					Run(gomock.Any(), "git", "log", "--format=%h|%s|%an|%ar", "v1.0.0..v1.2.0").
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
					Run(gomock.Any(), "git", "log", "--format=%h|%s|%an|%ar", "v1.0.0..v1.2.0").
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
		{
			name: "No commits found",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "log", "--format=%h|%s|%an|%ar", "v1.0.0..v1.2.0").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "")
							cmd.Stdout = bytes.NewBufferString("")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.ErrNoGitCommits,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			output, err := gitClient.Log(context.Background(), []string{"v1.0.0..v1.2.0"})

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, output)
			}
		})
	}
}

func Test_Status(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockExecutor := NewMockExecutor(ctrl)
	log := newStubLogger()

	gitClient := &client{
		executor: mockExecutor,
		log:      log,
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
					Run(gomock.Any(), "git", "diff", "--staged", "--name-status").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "A\tCLAUDE.md")
							cmd.Stdout = bytes.NewBufferString("A\tCLAUDE.md")
							return cmd
						})
			},
			expected: result{
				output: "A\tCLAUDE.md",
				err:    nil,
			},
		},
		{
			name: "Failure",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--name-status").
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
		{
			name: "No changes to commit",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--name-status").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "")
							cmd.Stdout = bytes.NewBufferString("")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.ErrNoGitChanges,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			output, err := gitClient.Status(ctx)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, output)
			}
		})
	}
}

func Test_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	log := newStubLogger()

	gitClient := &client{
		executor: mockExecutor,
		log:      log,
	}

	type result struct {
		output string
		err    error
	}

	tests := []struct {
		name     string
		message  string
		before   func()
		expected result
	}{
		{
			name:    "Success",
			message: "test commit",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "commit", "-m", "test commit").
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("echo", "commit successful")
							cmd.Stdout = bytes.NewBufferString("commit successful\n")
							return cmd
						})
			},
			expected: result{
				output: "commit successful\n",
				err:    nil,
			},
		},
		{
			name:    "Failure",
			message: "test commit",
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
		{
			name:    "Empty commit message",
			message: "",
			before:  func() {},
			expected: result{
				output: "",
				err:    errors.ErrCommitMessageEmpty,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			output, err := gitClient.Commit(context.Background(), tt.message)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.output, output)
			}
		})
	}
}
