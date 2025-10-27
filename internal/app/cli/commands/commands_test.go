package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

func Test_ProvideCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	params := CommandsParams{
		GitClient: mockGit,
		GPTClient: mockGPT,
		Log:       mockLogger,
		Spinner:   spinner.NewSpinner,
	}

	result := provideCommands(params)

	assert.NotNil(t, result.Help)
	assert.NotNil(t, result.Version)
	assert.NotNil(t, result.Changelog)
	assert.NotNil(t, result.Commit)
}
