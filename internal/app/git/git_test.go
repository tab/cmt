package git

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_NewGitClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.DefaultConfig()
	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	gitClient := NewGitClient(cfg, mockExecutor, mockLogger)
	assert.NotNil(t, gitClient)

	instance, ok := gitClient.(*client)
	assert.True(t, ok)
	assert.NotNil(t, instance.cfg)
}

func Test_Diff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := config.DefaultConfig()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	gitClient := &client{
		cfg:      cfg,
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
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
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
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()
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
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	cfg := config.DefaultConfig()

	gitClient := &client{
		cfg:      cfg,
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
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
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
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()
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

func Test_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	cfg := config.DefaultConfig()
	cfg.Editor = "test-editor"

	gitClient := &client{
		cfg:      cfg,
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
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
				mockExecutor.EXPECT().
					Run(gomock.Any(), "test-editor", gomock.Any()).
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							err := os.WriteFile(args[0], []byte("edited commit message"), 0644)
							if err != nil {
								t.Fatalf("Failed to write to temp file: %v", err)
							}
							cmd := exec.Command("echo", "")
							return cmd
						})
			},
			expected: result{
				output: "edited commit message",
				err:    nil,
			},
		},
		{
			name: "Failure to run editor",
			before: func() {
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()
				mockExecutor.EXPECT().
					Run(gomock.Any(), "test-editor", gomock.Any()).
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							cmd := exec.Command("false")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.New("error running editor"),
			},
		},
		{
			name: "Failure to read edited file",
			before: func() {
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()
				mockExecutor.EXPECT().
					Run(gomock.Any(), "test-editor", gomock.Any()).
					Return(nil).
					DoAndReturn(
						func(ctx context.Context, name string, args ...string) *exec.Cmd {
							os.Remove(args[0])
							cmd := exec.Command("echo", "")
							return cmd
						})
			},
			expected: result{
				output: "",
				err:    errors.New("failed to read file"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			output, err := gitClient.Edit(context.Background(), "original message")

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
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	cfg := config.DefaultConfig()

	gitClient := &client{
		cfg:      cfg,
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
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
				mockLogger.EXPECT().Info().Return(mockEvent).AnyTimes()
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
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()
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

			output, err := gitClient.Commit(context.Background(), "test commit")

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
