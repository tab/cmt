package model

import (
	"strings"

	"cmt/internal/app/cli/components"
	"cmt/internal/config/logger"
)

// FormatLogEntry formats a log entry with appropriate color based on level
func FormatLogEntry(entry *logger.LogEntry) string {
	formatted := entry.Format()

	level := strings.ToLower(entry.Level)
	switch level {
	case "debug":
		return components.LogDebugStyle.Render(formatted)
	case "info":
		return components.LogInfoStyle.Render(formatted)
	case "warn", "warning":
		return components.LogWarnStyle.Render(formatted)
	case "error":
		return components.LogErrorStyle.Render(formatted)
	case "trace":
		return components.LogTraceStyle.Render(formatted)
	default:
		return formatted
	}
}
