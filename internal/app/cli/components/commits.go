package components

import (
	"strings"
)

// ConvertCommitLogForGPT converts the detailed log format back to simple format for GPT
// Input format: hash|subject|author|date
// Output format: hash subject
func ConvertCommitLogForGPT(log string) string {
	if log == "" {
		return ""
	}

	var b strings.Builder
	lines := strings.Split(strings.TrimSpace(log), "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 2 {
			continue
		}

		b.WriteString(parts[0])
		b.WriteString(" ")
		b.WriteString(parts[1])
		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}
