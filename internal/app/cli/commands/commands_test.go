package commands

import (
	"testing"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_provideCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	log := stubLogger{log: zerolog.Nop()}

	params := CommandsParams{
		GitClient:      mockGit,
		GPTClient:      mockGPT,
		Log:            log,
		SpinnerFactory: spinner.NewSpinner,
	}

	result := provideCommands(params)

	assert.NotNil(t, result.Help)
	assert.NotNil(t, result.Version)
	assert.NotNil(t, result.Changelog)
	assert.NotNil(t, result.Commit)
}
