package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/config"
)

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedPrefix string
	}{
		{name: "no arguments - default to commit", args: []string{}, expectedPrefix: ""},
		{name: "prefix command with value", args: []string{"prefix", "JIRA-123"}, expectedPrefix: "JIRA-123"},
		{name: "prefix command without value", args: []string{"prefix"}, expectedPrefix: ""},
		{name: "prefix short flag with value", args: []string{"-p", "JIRA-456"}, expectedPrefix: "JIRA-456"},
		{name: "prefix long flag with value", args: []string{"--prefix", "JIRA-789"}, expectedPrefix: "JIRA-789"},
		{name: "single argument as prefix", args: []string{"CUSTOM-PREFIX"}, expectedPrefix: "CUSTOM-PREFIX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockWorkflow := NewMockworkflowService(ctrl)

			m := New(context.Background(), &config.Config{}, mockWorkflow, newTestLogger(), tt.args)

			assert.Equal(t, FlowCommit, m.UserFlow)
			assert.Equal(t, tt.expectedPrefix, m.Prefix)
		})
	}
}
