package model

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"cmt/internal/config/logger"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiRegexp.ReplaceAllString(input, "")
}

func Test_FormatLogEntry(t *testing.T) {
	baseTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		level string
	}{
		{
			name:  "debug level",
			level: "debug",
		},
		{
			name:  "info level",
			level: "INFO",
		},
		{
			name:  "warn level",
			level: "warn",
		},
		{
			name:  "warning alias",
			level: "warning",
		},
		{
			name:  "error level",
			level: "ERROR",
		},
		{
			name:  "trace level",
			level: "trace",
		},
		{
			name:  "unknown level",
			level: "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &logger.LogEntry{
				Timestamp: baseTime,
				Level:     tt.level,
				Message:   "log message",
				Fields: map[string]interface{}{
					"request_id": 42,
				},
			}

			formatted := entry.Format()
			result := FormatLogEntry(entry)
			assert.Equal(t, formatted, stripANSI(result))
		})
	}
}
