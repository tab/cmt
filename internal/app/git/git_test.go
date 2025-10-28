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

// fakeCommandWithOutput returns a function that creates an exec.Cmd that outputs the given string.
func fakeCommandWithOutput(output string) func(context.Context, string, ...string) *exec.Cmd {
	return func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.Command("printf", "%s", output)
		cmd.Stdout = bytes.NewBufferString(output)
		return cmd
	}
}

// fakeFailingCommand returns a function that creates an exec.Cmd that fails.
func fakeFailingCommand() func(context.Context, string, ...string) *exec.Cmd {
	return func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
}

// fakeEmptyCommand returns a function that creates an exec.Cmd that outputs nothing.
func fakeEmptyCommand() func(context.Context, string, ...string) *exec.Cmd {
	return func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.Command("printf", "")
		cmd.Stdout = bytes.NewBufferString("")
		return cmd
	}
}

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_NewGitClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	gitClient := NewGitClient(mockExecutor, mockLogger)
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
	nopLogger := zerolog.Nop()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	// Allow any logger calls since we're testing git functionality, not logging
	mockLogger.EXPECT().Debug().AnyTimes().Return(nopLogger.Debug())
	mockLogger.EXPECT().Error().AnyTimes().Return(nopLogger.Error())
	mockLogger.EXPECT().Info().AnyTimes().Return(nopLogger.Info())

	gitClient := &client{
		executor: mockExecutor,
		log:      mockLogger,
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
					DoAndReturn(fakeCommandWithOutput("mock diff output"))
			},
			expected: result{
				output: "mock diff output",
				err:    nil,
			},
		},
		{
			name: "Failure when diff command fails",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").
					Return(nil).
					DoAndReturn(fakeFailingCommand())
			},
			expected: result{
				output: "",
				err:    errors.New("failed to load git diff"),
			},
		},
		{
			name: "Failure with no changes",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines").
					Return(nil).
					DoAndReturn(fakeEmptyCommand())
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

//nolint:dupl
func Test_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nopLogger := zerolog.Nop()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	// Allow any logger calls since we're testing git functionality, not logging
	mockLogger.EXPECT().Debug().AnyTimes().Return(nopLogger.Debug())
	mockLogger.EXPECT().Error().AnyTimes().Return(nopLogger.Error())
	mockLogger.EXPECT().Info().AnyTimes().Return(nopLogger.Info())

	gitClient := &client{
		executor: mockExecutor,
		log:      mockLogger,
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
					DoAndReturn(fakeCommandWithOutput("mock log output"))
			},
			expected: result{
				output: "mock log output",
				err:    nil,
			},
		},
		{
			name: "Failure when log command fails",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "log", "--format=%h|%s|%an|%ar", "v1.0.0..v1.2.0").
					Return(nil).
					DoAndReturn(fakeFailingCommand())
			},
			expected: result{
				output: "",
				err:    errors.New("failed to load git log"),
			},
		},
		{
			name: "Failure with no commits",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "log", "--format=%h|%s|%an|%ar", "v1.0.0..v1.2.0").
					Return(nil).
					DoAndReturn(fakeEmptyCommand())
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

//nolint:dupl
func Test_Status(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	nopLogger := zerolog.Nop()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	// Allow any logger calls since we're testing git functionality, not logging
	mockLogger.EXPECT().Debug().AnyTimes().Return(nopLogger.Debug())
	mockLogger.EXPECT().Error().AnyTimes().Return(nopLogger.Error())
	mockLogger.EXPECT().Info().AnyTimes().Return(nopLogger.Info())

	gitClient := &client{
		executor: mockExecutor,
		log:      mockLogger,
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
					DoAndReturn(fakeCommandWithOutput("A\tCLAUDE.md"))
			},
			expected: result{
				output: "A\tCLAUDE.md",
				err:    nil,
			},
		},
		{
			name: "Failure when status command fails",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--name-status").
					Return(nil).
					DoAndReturn(fakeFailingCommand())
			},
			expected: result{
				output: "",
				err:    errors.New("failed to load git diff"),
			},
		},
		{
			name: "Failure with no changes",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "diff", "--staged", "--name-status").
					Return(nil).
					DoAndReturn(fakeEmptyCommand())
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

	nopLogger := zerolog.Nop()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	// Allow any logger calls since we're testing git functionality, not logging
	mockLogger.EXPECT().Debug().AnyTimes().Return(nopLogger.Debug())
	mockLogger.EXPECT().Error().AnyTimes().Return(nopLogger.Error())
	mockLogger.EXPECT().Info().AnyTimes().Return(nopLogger.Info())

	gitClient := &client{
		executor: mockExecutor,
		log:      mockLogger,
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
					DoAndReturn(fakeCommandWithOutput("commit successful\n"))
			},
			expected: result{
				output: "commit successful\n",
				err:    nil,
			},
		},
		{
			name:    "Failure when commit command fails",
			message: "test commit",
			before: func() {
				mockExecutor.EXPECT().
					Run(gomock.Any(), "git", "commit", "-m", "test commit").
					Return(nil).
					DoAndReturn(fakeFailingCommand())
			},
			expected: result{
				output: "",
				err:    errors.New("failed to commit changes"),
			},
		},
		{
			name:    "Failure with empty message",
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
