package workflow

import (
	"context"
	stdErrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func newTestLogger() logger.Logger {
	cfg := config.DefaultConfig()
	cfg.Logging.Format = logger.JSONFormat
	cfg.Logging.Level = logger.DebugLevel
	return logger.NewLogger(cfg)
}

func Test_GenerateCommit(t *testing.T) {
	t.Run("success with prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Diff(gomock.Any()).Return("diff content", nil)
		mockGit.EXPECT().Status(gomock.Any()).Return("A\tfile.go", nil)
		mockGpt.EXPECT().FetchCommitMessage(gomock.Any(), "diff content").Return("feat: add feature", nil)

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateCommit(context.Background(), "JIRA-123")
		assert.NoError(t, err)
		assert.Equal(t, "JIRA-123 feat: add feature", result.Message)
		assert.NotNil(t, result.FileTree)
	})

	t.Run("status error tolerated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Diff(gomock.Any()).Return("diff content", nil)
		mockGit.EXPECT().Status(gomock.Any()).Return("", stdErrors.New("status error"))
		mockGpt.EXPECT().FetchCommitMessage(gomock.Any(), "diff content").Return("fix: resolve issue", nil)

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateCommit(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, "fix: resolve issue", result.Message)
		assert.Nil(t, result.FileTree)
	})

	t.Run("diff failure propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Diff(gomock.Any()).Return("", stdErrors.New("diff error"))

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateCommit(context.Background(), "")
		assert.Error(t, err)
		assert.Empty(t, result.Message)
	})

	t.Run("gpt failure propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Diff(gomock.Any()).Return("diff content", nil)
		mockGit.EXPECT().Status(gomock.Any()).Return("A\tfile.go", nil)
		mockGpt.EXPECT().FetchCommitMessage(gomock.Any(), "diff content").Return("", stdErrors.New("gpt error"))

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateCommit(context.Background(), "")
		assert.Error(t, err)
		assert.Empty(t, result.Message)
	})
}

func Test_GenerateChangelog(t *testing.T) {
	t.Run("success without range", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Log(gomock.Any(), []string{}).Return("abc123|feat: add api|alice|2 days ago", nil)
		mockGpt.EXPECT().FetchChangelog(gomock.Any(), "abc123 feat: add api").Return("# CHANGELOG", nil)

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateChangelog(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, "# CHANGELOG", result.Content)
	})

	t.Run("success with range", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		between := "v1.0.0..v1.1.0"
		mockGit.EXPECT().Log(gomock.Any(), []string{between}).Return("def456|fix: bug|bob|yesterday", nil)
		mockGpt.EXPECT().FetchChangelog(gomock.Any(), "def456 fix: bug").Return("Changelog", nil)

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateChangelog(context.Background(), between)
		assert.NoError(t, err)
		assert.Equal(t, "Changelog", result.Content)
	})

	t.Run("git error propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Log(gomock.Any(), []string{}).Return("", stdErrors.New("git error"))

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateChangelog(context.Background(), "")
		assert.Error(t, err)
		assert.Empty(t, result.Content)
	})

	t.Run("no commits error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Log(gomock.Any(), []string{}).Return("", nil)

		service := NewService(mockGit, mockGpt, log)

		_, err := service.GenerateChangelog(context.Background(), "")
		assert.ErrorIs(t, err, errors.ErrNoGitCommits)
	})

	t.Run("gpt error propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGit := git.NewMockClient(ctrl)
		mockGpt := gpt.NewMockClient(ctrl)
		log := newTestLogger()

		mockGit.EXPECT().Log(gomock.Any(), []string{}).Return("abc123|feat: add api|alice|today", nil)
		mockGpt.EXPECT().FetchChangelog(gomock.Any(), "abc123 feat: add api").Return("", stdErrors.New("gpt error"))

		service := NewService(mockGit, mockGpt, log)

		result, err := service.GenerateChangelog(context.Background(), "")
		assert.Error(t, err)
		assert.Empty(t, result.Content)
	})
}
