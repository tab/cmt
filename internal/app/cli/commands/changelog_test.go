package commands

import (
	"context"
	"errors"
	"testing"

	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func Test_NewChangelogCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	log := newStubLogger()

	cmd := NewChangelogCommand(mockGit, mockGPT, log)
	assert.NotNil(t, cmd)
}

func Test_changelogCmd_Run(t *testing.T) {
	log := newStubLogger()

	tests := []struct {
		name           string
		args           []string
		before         func(*git.MockClient, *gpt.MockClient)
		expectedReturn int
	}{
		{
			name: "successful changelog generation with no args",
			args: []string{},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commit1\ncommit2", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commit1\ncommit2").
					Return("# Changelog\n\n- Feature 1\n- Feature 2", nil)
			},
			expectedReturn: 0,
		},
		{
			name: "successful changelog generation with range",
			args: []string{"v1.0..v2.0"},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commit1\ncommit2", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commit1\ncommit2").
					Return("# Changelog\n\n- Feature 1", nil)
			},
			expectedReturn: 0,
		},
		{
			name: "successful changelog with multiple args",
			args: []string{"HEAD~10", "HEAD"},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commits", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commits").
					Return("# Changelog", nil)
			},
			expectedReturn: 0,
		},
		{
			name: "args with whitespace are trimmed",
			args: []string{"  v1.0  ", "", "  v2.0  "},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commits", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commits").
					Return("# Changelog", nil)
			},
			expectedReturn: 0,
		},
		{
			name: "git log fails",
			args: []string{},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("", errors.New("git error"))
			},
			expectedReturn: 1,
		},
		{
			name: "gpt fetch fails",
			args: []string{},
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().
					Log(gomock.Any(), gomock.Any()).
					Return("commits", nil)

				mockGPT.EXPECT().
					FetchChangelog(gomock.Any(), "commits").
					Return("", errors.New("api error"))
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

			tt.before(mockGit, mockGPT)

			cmd := NewChangelogCommand(mockGit, mockGPT, log)
			result := cmd.Run(context.Background(), tt.args)

			assert.Equal(t, tt.expectedReturn, result)
		})
	}
}
