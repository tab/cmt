package components

import (
	"strings"

	"cmt/internal/config"
)

// RenderHeader renders the application title header
func RenderHeader(title string, width int) string {
	return SectionHeader.Width(width).Render(">_ " + title)
}

// RenderAppHeader renders the main app header with app name
func RenderAppHeader(width int) string {
	return RenderHeader(config.AppName, width)
}

// RenderHints renders keyboard hints at the bottom
func RenderHints(hints []string) string {
	hintsText := strings.Join(hints, " • ")
	return HelpText.Render(hintsText)
}

// RenderCommitHints renders hints for commit review state
func RenderCommitHints(hasLogs bool, width int) string {
	hints := []string{
		KeyHint.Render("[a]") + " accept",
		KeyHint.Render("[e]") + " edit",
		KeyHint.Render("[r]") + " refresh",
	}

	if hasLogs {
		hints = append(hints, KeyHint.Render("[l]")+" logs")
	}

	hints = append(hints, KeyHint.Render("[q]")+" quit")

	return RenderHints(hints)
}

// RenderPanel renders a bordered panel with content
func RenderPanel(content string, width, height int) string {
	return PanelStyle(width, height).Render(content)
}

// RenderError renders an error message
func RenderError(err error) string {
	if err == nil {
		return ""
	}
	return StatusError.Render("❌ " + err.Error())
}

// RenderLogsHints renders hints for logs view
func RenderLogsHints() string {
	hints := []string{
		KeyHint.Render("[q]") + " close",
		KeyHint.Render("[↑/↓]") + " scroll",
	}

	return RenderHints(hints)
}
