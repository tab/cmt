package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/commands"
	"cmt/internal/app/errors"
)

func Test_NewRunner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := Params{
		Help:      commands.NewMockCommand(ctrl),
		Version:   commands.NewMockCommand(ctrl),
		Changelog: commands.NewMockCommand(ctrl),
		Commit:    commands.NewMockCommand(ctrl),
	}

	instance := NewRunner(params)
	assert.NotNil(t, instance)
}

func Test_Resolve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	helpCmd := commands.NewMockCommand(ctrl)
	versionCmd := commands.NewMockCommand(ctrl)
	changelogCmd := commands.NewMockCommand(ctrl)
	commitCmd := commands.NewMockCommand(ctrl)

	params := Params{
		Help:      helpCmd,
		Version:   versionCmd,
		Changelog: changelogCmd,
		Commit:    commitCmd,
	}

	instance := NewRunner(params)

	tests := []struct {
		name          string
		args          []string
		expectedCmd   commands.Command
		expectedArgs  []string
		expectedError error
	}{
		{
			name:          "Success without args",
			args:          []string{},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with help command",
			args:          []string{"help"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with help flag long",
			args:          []string{"--help"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with help flag short",
			args:          []string{"-h"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with version command",
			args:          []string{"version"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with version flag long",
			args:          []string{"--version"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with version flag short",
			args:          []string{"-v"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with changelog command",
			args:          []string{"changelog"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with changelog flag long",
			args:          []string{"--changelog"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with changelog flag short",
			args:          []string{"-c"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with changelog range",
			args:          []string{"changelog", "v1.0..v2.0"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"v1.0..v2.0"},
			expectedError: nil,
		},
		{
			name:          "Success with changelog range boundaries",
			args:          []string{"changelog", "HEAD~10", "HEAD"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"HEAD~10", "HEAD"},
			expectedError: nil,
		},
		{
			name:          "Success with changelog flag long range",
			args:          []string{"--changelog", "v1.0..v2.0"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"v1.0..v2.0"},
			expectedError: nil,
		},
		{
			name:          "Success with changelog flag short range",
			args:          []string{"-c", "v1.0..v2.0"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"v1.0..v2.0"},
			expectedError: nil,
		},
		{
			name:          "Success with prefix flag long",
			args:          []string{"--prefix", "feat:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"--prefix", "feat:"},
			expectedError: nil,
		},
		{
			name:          "Success with prefix flag short",
			args:          []string{"-p", "fix:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"-p", "fix:"},
			expectedError: nil,
		},
		{
			name:          "Success with prefix command",
			args:          []string{"prefix", "chore:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"prefix", "chore:"},
			expectedError: nil,
		},
		{
			name:          "Success with prefix inline",
			args:          []string{"--prefix=docs:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"--prefix=docs:"},
			expectedError: nil,
		},
		{
			name:          "Success with case insensitive help",
			args:          []string{"HELP"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with case insensitive version",
			args:          []string{"VERSION"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Success with case insensitive changelog",
			args:          []string{"CHANGELOG"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "Failure with unknown command",
			args:          []string{"unknown"},
			expectedCmd:   nil,
			expectedArgs:  nil,
			expectedError: errors.ErrUnknownCommand,
		},
		{
			name:          "Failure with unknown flag",
			args:          []string{"--unknown"},
			expectedCmd:   nil,
			expectedArgs:  nil,
			expectedError: errors.ErrUnknownCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args, err := instance.Resolve(tt.args)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, cmd)
				assert.Nil(t, args)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmd, cmd)
				assert.Equal(t, tt.expectedArgs, args)
			}
		})
	}
}
