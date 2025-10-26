package commands

import (
	"context"
	"strings"
	"testing"

	"cmt/internal/config"
	"github.com/stretchr/testify/assert"
)

func Test_NewHelpCommand(t *testing.T) {
	cmd := NewHelpCommand()
	assert.NotNil(t, cmd)
}

func Test_helpCmd_Run(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedReturn int
	}{
		{
			name:           "run help command with no args",
			args:           []string{},
			expectedReturn: 0,
		},
		{
			name:           "run help command with args",
			args:           []string{"some", "args"},
			expectedReturn: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewHelpCommand()
			result := cmd.Run(context.Background(), tt.args)
			assert.Equal(t, tt.expectedReturn, result)
		})
	}
}

func Test_GetUsage(t *testing.T) {
	usage := GetUsage()

	tests := []struct {
		name     string
		contains string
	}{
		{
			name:     "contains app name",
			contains: config.AppName,
		},
		{
			name:     "contains version",
			contains: config.Version,
		},
		{
			name:     "contains description",
			contains: config.AppDescription,
		},
		{
			name:     "contains usage section",
			contains: "Usage:",
		},
		{
			name:     "contains commands section",
			contains: "Commands:",
		},
		{
			name:     "contains examples section",
			contains: "Examples:",
		},
		{
			name:     "contains navigation section",
			contains: "Navigation:",
		},
		{
			name:     "contains environment section",
			contains: "Environment:",
		},
		{
			name:     "contains changelog command",
			contains: "changelog",
		},
		{
			name:     "contains version command",
			contains: "version",
		},
		{
			name:     "contains help command",
			contains: "help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, strings.Contains(usage, tt.contains),
				"expected usage to contain %q", tt.contains)
		})
	}
}
