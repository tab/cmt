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

func Test_NewChangelog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.DefaultConfig()
	mockClient := git.NewMockClient(ctrl)
	mockModel := gpt.NewMockClient(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	mockLoader := utils.NewMockLoader(ctrl)

	cmd := NewChangelog(cfg, mockClient, mockModel, mockLoader, mockLogger)
	assert.NotNil(t, cmd)

	instance, ok := cmd.(*changelog)
	assert.True(t, ok)
	assert.NotNil(t, instance)
}

func Test_Changelog_Generate(t *testing.T) {
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

	cmd := NewChangelog(cfg, mockClient, mockModel, mockLoader, mockLogger)

	type result struct {
		output string
		err    error
	}

	tests := []struct {
		name   string
		args   []string
		before func()
		expect result
	}{
		{
			name: "Success",
			args: []string{"2606b09..5e3ac73"},
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Log(gomock.Any(), []string{"2606b09..5e3ac73"}).Return("git log content", nil)
				mockModel.EXPECT().FetchChangelog(gomock.Any(), "git log content").Return("content", nil)
			},
			expect: result{
				output: "💬 Changelog: \n\ncontent\n",
				err:    nil,
			},
		},
		{
			name: "Success with tags",
			args: []string{"v1.0.0..v1.1.0"},
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Log(gomock.Any(), []string{"v1.0.0..v1.1.0"}).Return("git log content", nil)
				mockModel.EXPECT().FetchChangelog(gomock.Any(), "git log content").Return("content", nil)
			},
			expect: result{
				output: "💬 Changelog: \n\ncontent\n",
				err:    nil,
			},
		},
		{
			name: "No arguments",
			args: []string{},
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Log(gomock.Any(), []string{}).Return("git log content", nil)
				mockModel.EXPECT().FetchChangelog(gomock.Any(), "git log content").Return("", errors.ErrFailedToLoadGitLog)
			},
			expect: result{
				output: "❌ Error getting model response: failed to load git log\n",
				err:    errors.ErrFailedToLoadGitLog,
			},
		},
		{
			name: "Failed to fetch git log",
			args: []string{"2606b09..5e3ac73"},
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Log(gomock.Any(), []string{"2606b09..5e3ac73"}).Return("", errors.ErrNoGitChanges)
			},
			expect: result{
				output: "⚠️ No changes found in the git log\n",
				err:    errors.ErrNoGitChanges,
			},
		},
		{
			name: "Failed to generate changelog",
			args: []string{"2606b09..5e3ac73"},
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Log(gomock.Any(), []string{"2606b09..5e3ac73"}).Return("git log content", nil)
				mockModel.EXPECT().FetchChangelog(gomock.Any(), "git log content").Return("", errors.ErrNoResponse)
			},
			expect: result{
				output: "⚠️ No response from GPT\n",
				err:    errors.ErrNoResponse,
			},
		},
		{
			name: "Failed to parse JSON response",
			args: []string{"2606b09..5e3ac73"},
			before: func() {
				mockLoader.EXPECT().Start().Times(1)
				mockLoader.EXPECT().Stop().Times(1)

				mockLogger.EXPECT().Debug().Return(mockDebugEvent).AnyTimes()

				mockClient.EXPECT().Log(gomock.Any(), []string{"2606b09..5e3ac73"}).Return("git log content", nil)
				mockModel.EXPECT().FetchChangelog(gomock.Any(), "git log content").Return("", errors.ErrFailedToParseJSON)
			},
			expect: result{
				output: "⚠️ Failed to parse JSON response\n",
				err:    errors.ErrFailedToParseJSON,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			defer func() {
				os.Stdout = oldStdout
			}()

			err := cmd.Generate(ctx, tt.args)

			w.Close()
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			if tt.expect.err != nil {
				assert.Equal(t, tt.expect.err, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expect.output, output)
		})
	}
}
