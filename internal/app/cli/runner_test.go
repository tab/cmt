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
			name:          "no arguments returns commit command",
			args:          []string{},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "help command",
			args:          []string{"help"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "--help flag",
			args:          []string{"--help"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "-h flag",
			args:          []string{"-h"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "version command",
			args:          []string{"version"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "--version flag",
			args:          []string{"--version"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "-v flag",
			args:          []string{"-v"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "changelog command",
			args:          []string{"changelog"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "--changelog flag",
			args:          []string{"--changelog"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "-c flag",
			args:          []string{"-c"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "changelog with range",
			args:          []string{"changelog", "v1.0..v2.0"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"v1.0..v2.0"},
			expectedError: nil,
		},
		{
			name:          "changelog with multiple args",
			args:          []string{"changelog", "HEAD~10", "HEAD"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"HEAD~10", "HEAD"},
			expectedError: nil,
		},
		{
			name:          "--changelog with range",
			args:          []string{"--changelog", "v1.0..v2.0"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"v1.0..v2.0"},
			expectedError: nil,
		},
		{
			name:          "-c with range",
			args:          []string{"-c", "v1.0..v2.0"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{"v1.0..v2.0"},
			expectedError: nil,
		},
		{
			name:          "--prefix with value",
			args:          []string{"--prefix", "feat:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"--prefix", "feat:"},
			expectedError: nil,
		},
		{
			name:          "-p with value",
			args:          []string{"-p", "fix:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"-p", "fix:"},
			expectedError: nil,
		},
		{
			name:          "prefix with value",
			args:          []string{"prefix", "chore:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"prefix", "chore:"},
			expectedError: nil,
		},
		{
			name:          "--prefix=value",
			args:          []string{"--prefix=docs:"},
			expectedCmd:   commitCmd,
			expectedArgs:  []string{"--prefix=docs:"},
			expectedError: nil,
		},
		{
			name:          "case insensitive help",
			args:          []string{"HELP"},
			expectedCmd:   helpCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "case insensitive version",
			args:          []string{"VERSION"},
			expectedCmd:   versionCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "case insensitive changelog",
			args:          []string{"CHANGELOG"},
			expectedCmd:   changelogCmd,
			expectedArgs:  []string{},
			expectedError: nil,
		},
		{
			name:          "unknown command",
			args:          []string{"unknown"},
			expectedCmd:   nil,
			expectedArgs:  nil,
			expectedError: errors.ErrUnknownCommand,
		},
		{
			name:          "unknown flag",
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
