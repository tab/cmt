package components

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"cmt/internal/config"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiRegexp.ReplaceAllString(input, "")
}

func Test_RenderHeader(t *testing.T) {
	result := stripANSI(RenderHeader("Test Title", 30))
	assert.Contains(t, result, ">_ Test Title")
}

func Test_RenderAppHeader(t *testing.T) {
	result := stripANSI(RenderAppHeader(40))
	assert.Contains(t, result, config.AppName)
}

func Test_RenderHints(t *testing.T) {
	result := stripANSI(RenderHints([]string{"hint one", "hint two"}))
	clean := strings.ReplaceAll(result, "  ", " ")
	assert.Contains(t, clean, "hint one")
	assert.Contains(t, clean, "hint two")
}

func Test_RenderCommitHints(t *testing.T) {
	tests := []struct {
		name        string
		hasLogs     bool
		expectLogs  bool
		expectCount []string
	}{
		{
			name:       "without logs",
			hasLogs:    false,
			expectLogs: false,
			expectCount: []string{
				"accept",
				"edit",
				"refresh",
				"quit",
			},
		},
		{
			name:       "with logs",
			hasLogs:    true,
			expectLogs: true,
			expectCount: []string{
				"accept",
				"edit",
				"refresh",
				"logs",
				"quit",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripANSI(RenderCommitHints(tt.hasLogs, 60))
			if tt.expectLogs {
				assert.Contains(t, result, "logs")
			} else {
				assert.NotContains(t, result, "logs")
			}

			for _, expected := range tt.expectCount {
				assert.Contains(t, result, expected)
			}
		})
	}
}

func Test_RenderPanel(t *testing.T) {
	content := "panel content"
	result := stripANSI(RenderPanel(content, 20, 5))
	assert.Contains(t, result, content)
}

func Test_RenderError(t *testing.T) {
	assert.Equal(t, "", RenderError(nil))

	err := errors.New("failed to load")
	result := stripANSI(RenderError(err))

	assert.Contains(t, result, "failed to load")
	assert.Contains(t, result, "❌")
}

func Test_RenderLogsHints(t *testing.T) {
	result := stripANSI(RenderLogsHints())
	assert.Contains(t, result, "close")
	assert.Contains(t, result, "scroll")
}
