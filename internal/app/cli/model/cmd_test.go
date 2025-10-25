package model

import (
	"context"
	stdErrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/workflow"
	"cmt/internal/config/logger"
)

func Test_Start(t *testing.T) {
	cases := []struct {
		name       string
		setupMocks func(ctrl *gomock.Controller) (workflowService, logger.Logger)
		model      Model
		assertFunc func(t *testing.T, msg teaMsg)
	}{
		{
			name: "commit_flow_success",
			setupMocks: func(ctrl *gomock.Controller) (workflowService, logger.Logger) {
				mockWorkflow := NewMockworkflowService(ctrl)
				mockWorkflow.EXPECT().GenerateCommit(gomock.Any(), "prefix").
					Return(workflow.CommitResult{Message: "feat: add feature"}, nil)
				return mockWorkflow, newTestLogger()
			},
			model: Model{
				Ctx:      context.Background(),
				Prefix:   "prefix",
				UserFlow: FlowCommit,
			},
			assertFunc: func(t *testing.T, msg teaMsg) {
				result, ok := msg.(Result)
				assert.True(t, ok)
				assert.NoError(t, result.Err)
				assert.Equal(t, "feat: add feature", result.Content)
			},
		},
		{
			name: "commit_flow_error",
			setupMocks: func(ctrl *gomock.Controller) (workflowService, logger.Logger) {
				mockWorkflow := NewMockworkflowService(ctrl)
				mockWorkflow.EXPECT().GenerateCommit(gomock.Any(), "prefix").
					Return(workflow.CommitResult{}, stdErrors.New("boom"))
				return mockWorkflow, newTestLogger()
			},
			model: Model{
				Ctx:      context.Background(),
				Prefix:   "prefix",
				UserFlow: FlowCommit,
			},
			assertFunc: func(t *testing.T, msg teaMsg) {
				result, ok := msg.(Result)
				assert.True(t, ok)
				assert.Error(t, result.Err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWorkflow, log := tc.setupMocks(ctrl)
			tc.model.Workflow = mockWorkflow
			tc.model.Log = log

			msg := Start(tc.model)()
			tc.assertFunc(t, msg)
		})
	}
}

type teaMsg interface{}
