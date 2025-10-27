package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

func Test_NewChangelogCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	cmd := NewChangelogCommand(mockGit, mockGPT, mockLogger)
	assert.NotNil(t, cmd)
}

func Test_ChangelogCmd_Run(t *testing.T) {
	nopLogger := zerolog.Nop()

	tests := []struct {
		name           string
		args           []string
		before         func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger)
		expectedReturn int
	}{
		{
			name: "Success without args",
			args: []string{},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)

				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commit1\ncommit2", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commit1\ncommit2").
					Return("# Changelog\n\n- Feature 1\n- Feature 2", nil)

				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
			},
			expectedReturn: 0,
		},
		{
			name: "Success with range arg",
			args: []string{"v1.0..v2.0"},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)

				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commit1\ncommit2", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commit1\ncommit2").
					Return("# Changelog\n\n- Feature 1", nil)

				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
			},
			expectedReturn: 0,
		},
		{
			name: "Success with multiple args",
			args: []string{"HEAD~10", "HEAD"},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)

				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commits", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commits").
					Return("# Changelog", nil)

				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
			},
			expectedReturn: 0,
		},
		{
			name: "Success with whitespace args",
			args: []string{"  v1.0  ", "", "  v2.0  "},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)

				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commits", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commits").
					Return("# Changelog", nil)

				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
			},
			expectedReturn: 0,
		},
		{
			name: "Failure when git log fails",
			args: []string{},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)

				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("", errors.New("git error"))

				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)
			},
			expectedReturn: 1,
		},
		{
			name: "Failure when g p t fetch fails",
			args: []string{},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)

				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commits", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commits").
					Return("", errors.New("api error"))

				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)
			},
			expectedReturn: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGit := git.NewMockClient(ctrl)
			mockGPT := gpt.NewMockClient(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)

			cmd := NewChangelogCommand(mockGit, mockGPT, mockLogger)

			tt.before(mockGit, mockGPT, mockLogger)

			result := cmd.Run(context.Background(), tt.args)
			assert.Equal(t, tt.expectedReturn, result)
		})
	}
}
