package cli

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/model"
	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_NewUI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockWorkflow := model.NewMockworkflowService(ctrl)

	cfg := &config.Config{}
	cfg.Logging.Level = "error"
	cfg.Logging.Format = "console"
	log := logger.NewLogger(cfg)

	ctx := context.Background()

	uiInstance := NewUI(ctx, cfg, mockGit, mockWorkflow, log)
	assert.NotNil(t, uiInstance)
	assert.IsType(t, &ui{}, uiInstance)
}

func Test_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockWorkflow := model.NewMockworkflowService(ctrl)
	mockProgram := NewMockProgram(ctrl)

	cfg := &config.Config{}
	cfg.Logging.Level = "error"
	cfg.Logging.Format = "console"
	log := logger.NewLogger(cfg)

	tests := []struct {
		name        string
		args        []string
		before      func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful commit flow with accept",
			args: []string{},
			before: func() {
				mockProgram.EXPECT().Run().Return(model.Model{
					State:      model.StateExit,
					UserAction: model.ActionAccept,
					UserFlow:   model.FlowCommit,
					Content:    "feat: add new feature",
					Err:        nil,
				}, nil)
				mockGit.EXPECT().Commit(gomock.Any(), "feat: add new feature").Return("commit output", nil)
			},
			expectError: false,
		},
		{
			name: "commit flow with cancel",
			args: []string{},
			before: func() {
				mockProgram.EXPECT().Run().Return(model.Model{
					State:      model.StateExit,
					UserAction: model.ActionCancel,
					UserFlow:   model.FlowCommit,
					Content:    "feat: add new feature",
					Err:        nil,
				}, nil)
			},
			expectError: false,
		},
		{
			name: "error in TUI",
			args: []string{},
			before: func() {
				mockProgram.EXPECT().Run().Return(model.Model{
					State: model.StateExit,
					Err:   errors.ErrNoGitChanges,
				}, nil)
			},
			expectError: false,
		},
		{
			name: "program returns error",
			args: []string{},
			before: func() {
				mockProgram.EXPECT().Run().Return(nil, errors.New("program error"))
			},
			expectError: true,
			errorMsg:    "TUI error: program error",
		},
		{
			name: "commit fails after accept",
			args: []string{},
			before: func() {
				mockProgram.EXPECT().Run().Return(model.Model{
					State:      model.StateExit,
					UserAction: model.ActionAccept,
					UserFlow:   model.FlowCommit,
					Content:    "feat: add new feature",
					Err:        nil,
				}, nil)
				mockGit.EXPECT().Commit(gomock.Any(), "feat: add new feature").Return("", errors.New("commit failed"))
			},
			expectError: true,
			errorMsg:    "commit failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			ui := &ui{
				ctx: context.Background(),
				cfg: &config.Config{},
				builder: func(model tea.Model, opts ...tea.ProgramOption) Program {
					return mockProgram
				},
				git:      mockGit,
				workflow: mockWorkflow,
				log:      log,
			}

			err := ui.Run(tt.args)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
