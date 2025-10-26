package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()
	assert.NotNil(t, cmd)
}

func Test_versionCmd_Run(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedReturn int
	}{
		{
			name:           "run version command with no args",
			args:           []string{},
			expectedReturn: 0,
		},
		{
			name:           "run version command with args",
			args:           []string{"some", "args"},
			expectedReturn: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewVersionCommand()
			result := cmd.Run(context.Background(), tt.args)
			assert.Equal(t, tt.expectedReturn, result)
		})
	}
}
