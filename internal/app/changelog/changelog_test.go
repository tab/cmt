package changelog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/workflow"
	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

//go:generate mockgen -source=../cli/workflow/service.go -destination=workflow_mock_test.go -package=changelog

func TestGenerator_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.DefaultConfig()
	cfg.Logging.Level = "error"
	log := logger.NewLogger(cfg)

	tests := []struct {
		name    string
		between string
		setup   func(mockWorkflow *MockService)
		wantErr bool
	}{
		{
			name:    "successful generation",
			between: "v1.0.0..v2.0.0",
			setup: func(mockWorkflow *MockService) {
				mockWorkflow.EXPECT().
					GenerateChangelog(gomock.Any(), "v1.0.0..v2.0.0").
					Return(workflow.ChangelogResult{Content: "# CHANGELOG"}, nil)
			},
			wantErr: false,
		},
		{
			name:    "generation error",
			between: "",
			setup: func(mockWorkflow *MockService) {
				mockWorkflow.EXPECT().
					GenerateChangelog(gomock.Any(), "").
					Return(workflow.ChangelogResult{}, errors.ErrNoGitCommits)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWorkflow := NewMockService(ctrl)
			tt.setup(mockWorkflow)

			gen := NewGenerator(context.Background(), mockWorkflow, log)
			err := gen.Generate(tt.between)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
