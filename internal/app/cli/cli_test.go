package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/commands"
)

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_NewCLI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunner := NewMockRunner(ctrl)
	cli := NewCLI(mockRunner)
	assert.NotNil(t, cli)
	assert.NotNil(t, cli.runner)
}

func Test_CLI_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunner := NewMockRunner(ctrl)
	mockCmd := commands.NewMockCommand(ctrl)

	tests := []struct {
		name           string
		args           []string
		before         func(*MockRunner, *commands.MockCommand)
		expectedReturn int
	}{
		{
			name: "Successful no-arg command",
			args: []string{},
			before: func(mockRunner *MockRunner, mockCmd *commands.MockCommand) {
				mockRunner.EXPECT().Resolve([]string{}).Return(mockCmd, []string{}, nil)
				mockCmd.EXPECT().Run(gomock.Any(), []string{}).Return(0)
			},
			expectedReturn: 0,
		},
		{
			name: "Successful help command",
			args: []string{"help"},
			before: func(mockRunner *MockRunner, mockCmd *commands.MockCommand) {
				mockRunner.EXPECT().Resolve([]string{"help"}).Return(mockCmd, []string{}, nil)
				mockCmd.EXPECT().Run(gomock.Any(), []string{}).Return(0)
			},
			expectedReturn: 0,
		},
		{
			name: "Successful version command",
			args: []string{"version"},
			before: func(mockRunner *MockRunner, mockCmd *commands.MockCommand) {
				mockRunner.EXPECT().Resolve([]string{"version"}).Return(mockCmd, []string{}, nil)
				mockCmd.EXPECT().Run(gomock.Any(), []string{}).Return(0)
			},
			expectedReturn: 0,
		},
		{
			name: "Successful changelog command",
			args: []string{"changelog sha1..sha2"},
			before: func(mockRunner *MockRunner, mockCmd *commands.MockCommand) {
				mockRunner.EXPECT().Resolve([]string{"changelog sha1..sha2"}).Return(mockCmd, []string{}, nil)
				mockCmd.EXPECT().Run(gomock.Any(), []string{}).Return(0)
			},
			expectedReturn: 0,
		},
		{
			name: "Failure help command",
			args: []string{"help"},
			before: func(mockRunner *MockRunner, mockCmd *commands.MockCommand) {
				mockRunner.EXPECT().Resolve([]string{"help"}).Return(mockCmd, []string{}, nil)
				mockCmd.EXPECT().Run(gomock.Any(), []string{}).Return(1)
			},
			expectedReturn: 1,
		},
		{
			name: "Failure unknown command",
			args: []string{"unknown"},
			before: func(mockRunner *MockRunner, mockCmd *commands.MockCommand) {
				mockRunner.EXPECT().Resolve([]string{"unknown"}).Return(nil, nil, assert.AnError)
			},
			expectedReturn: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(mockRunner, mockCmd)

			cli := NewCLI(mockRunner)
			result := cli.Run(context.Background(), tt.args)
			assert.Equal(t, tt.expectedReturn, result)
		})
	}
}
