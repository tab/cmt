package commands

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/app/utils"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_NewCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.DefaultConfig()
	mockClient := git.NewMockClient(ctrl)
	mockModel := gpt.NewMockClient(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	mockLoader := utils.NewMockLoader(ctrl)

	cmd := NewCommit(cfg, mockClient, mockModel, mockLoader, mockLogger)
	assert.NotNil(t, cmd)

	instance, ok := cmd.(*commit)
	assert.True(t, ok)
	assert.NotNil(t, instance)
}

func Test_Commit_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	mockClient := git.NewMockClient(ctrl)
	mockModel := gpt.NewMockClient(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	mockLoader := utils.NewMockLoader(ctrl)

	nopLogger := zerolog.Nop()
	mockDebugEvent := nopLogger.Debug()

	cmd := NewCommit(cfg, mockClient, mockModel, mockLoader, mockLogger)

	type result struct {
		output string
		err    error
	}

	tests := []struct {
		name   string
		args   []string
		input  string
		before func()
		expect result
	}{
		{
			name:  "Success",
			args:  []string{},
			input: "y\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("commit message", nil)
				mockClient.EXPECT().Commit(gomock.Any(), "commit message").Return("", nil)
			},
			expect: result{
				output: "💬 Message: commit message\n\nAccept, edit, or cancel? (y/e/n): 🚀 Changes committed:\n\n",
				err:    nil,
			},
		},
		{
			name:  "Success with prefix",
			args:  []string{"prefix task-123"},
			input: "y\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("commit message", nil)
				mockClient.EXPECT().Commit(gomock.Any(), "prefix task-123 commit message").Return("", nil)
			},
			expect: result{
				output: "💬 Message: prefix task-123 commit message\n\nAccept, edit, or cancel? (y/e/n): 🚀 Changes committed:\n\n",
				err:    nil,
			},
		},
		{
			name:  "Success with edit",
			args:  []string{},
			input: "e\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("commit message", nil)
				mockClient.EXPECT().Edit(gomock.Any(), "commit message").Return("edit message", nil)
				mockClient.EXPECT().Commit(gomock.Any(), "edit message").Return("", nil)
			},
			expect: result{
				output: "💬 Message: commit message\n\nAccept, edit, or cancel? (y/e/n): \n🧑🏻\u200d💻 Commit message was changed successfully!\n🚀 Changes committed:\n\n",
				err:    nil,
			},
		},
		{
			name:  "Success with cancel",
			args:  []string{},
			input: "n\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("commit message", nil)
			},
			expect: result{
				output: "💬 Message: commit message\n\nAccept, edit, or cancel? (y/e/n): ❌ Commit aborted\n",
				err:    nil,
			},
		},
		{
			name:  "Error on git diff",
			args:  []string{},
			input: "y\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("", errors.ErrNoGitChanges)
			},
			expect: result{
				output: "⚠️ No changes to commit\n",
				err:    errors.ErrNoGitChanges,
			},
		},
		{
			name:  "Error on model fetch",
			args:  []string{},
			input: "y\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("", errors.ErrNoResponse)
			},
			expect: result{
				output: "⚠️ No response from GPT\n",
				err:    errors.ErrNoResponse,
			},
		},
		{
			name:  "Error on commit",
			args:  []string{},
			input: "y\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("commit message", nil)
				mockClient.EXPECT().Commit(gomock.Any(), "commit message").Return("", errors.ErrCommitMessageEmpty)
			},
			expect: result{
				output: "💬 Message: commit message\n\nAccept, edit, or cancel? (y/e/n): ⚠️ Commit message is empty, commit aborted\n",
				err:    errors.ErrCommitMessageEmpty,
			},
		},
		{
			name:  "Error on edit",
			args:  []string{},
			input: "e\n",
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Diff(ctx).Return("diff", nil)
				mockModel.EXPECT().FetchCommitMessage(ctx, "diff").Return("commit message", nil)
				mockClient.EXPECT().Edit(gomock.Any(), "commit message").Return("", errors.ErrFailedToRunEditor)
			},
			expect: result{
				output: "💬 Message: commit message\n\nAccept, edit, or cancel? (y/e/n): ❌ Error editing commit message: error running editor\n",
				err:    errors.ErrFailedToRunEditor,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			oldStdout := os.Stdout
			oldStdin := os.Stdin

			stdoutReader, stdoutWriter, err := os.Pipe()
			assert.NoError(t, err)

			stdinReader, stdinWriter, err := os.Pipe()
			assert.NoError(t, err)

			os.Stdout = stdoutWriter
			os.Stdin = stdinReader

			defer func() {
				os.Stdout = oldStdout
				os.Stdin = oldStdin
			}()

			_, err = stdinWriter.WriteString(tt.input)
			assert.NoError(t, err)
			stdinWriter.Close()

			err = cmd.Generate(ctx, tt.args)

			stdoutWriter.Close()
			var buf bytes.Buffer
			_, errCopy := io.Copy(&buf, stdoutReader)
			assert.NoError(t, errCopy)
			output := buf.String()

			if tt.expect.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expect.output, output)
		})
	}
}
