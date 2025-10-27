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

func Test_HelpCmd_Run(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedReturn int
	}{
		{
			name:           "Success without args",
			args:           []string{},
			expectedReturn: 0,
		},
		{
			name:           "Success with args",
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
			name:     "Success with app name",
			contains: config.AppName,
		},
		{
			name:     "Success with version",
			contains: config.Version,
		},
		{
			name:     "Success with description",
			contains: config.AppDescription,
		},
		{
			name:     "Success with usage section",
			contains: "Usage:",
		},
		{
			name:     "Success with commands section",
			contains: "Commands:",
		},
		{
			name:     "Success with examples section",
			contains: "Examples:",
		},
		{
			name:     "Success with navigation section",
			contains: "Navigation:",
		},
		{
			name:     "Success with environment section",
			contains: "Environment:",
		},
		{
			name:     "Success with changelog command",
			contains: "changelog",
		},
		{
			name:     "Success with version command",
			contains: "version",
		},
		{
			name:     "Success with help command",
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
