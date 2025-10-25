package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"cmt/internal/app/cli/components"
)

// View renders the current state (Bubble Tea lifecycle)
func (m Model) View() string {
	switch m.State {
	case StateInit, StateFetch:
		return m.viewFetch()

	case StateViewCommit:
		return m.viewReview()

	case StateEditCommit:
		return m.viewEdit()

	case StateViewLogs:
		return m.viewLogs()

	case StateExit:
		return ""

	default:
		return "Unknown state"
	}
}

// viewFetch renders the fetching state
func (m Model) viewFetch() string {
	header := components.RenderAppHeader(m.Width)
	loading := components.SpinnerStyle.Render(m.Spinner.View() + "loading...")

	return lipgloss.JoinVertical(lipgloss.Left, header, loading)
}

// viewReview renders the review state
func (m Model) viewReview() string {
	if m.Width < 40 || m.Height < 10 {
		return "Terminal too small"
	}

	leftWidth := int(float64(m.Width) * 0.28)
	rightWidth := m.Width - leftWidth
	height := m.Height - 6

	header := components.RenderAppHeader(m.Width)
	leftPanel := m.renderFileBrowser(leftWidth, height)
	rightPanel := m.renderCommitMessage(rightWidth, height)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	hints := components.RenderCommitHints(m.LogBuffer != nil, m.Width)

	return lipgloss.JoinVertical(lipgloss.Left, header, panels, hints)
}

// renderFileBrowser renders the file browser panel
func (m Model) renderFileBrowser(width, height int) string {
	if m.FileTree == nil {
		return components.RenderPanel(components.LabelMedium.Render("No staged files"), width, height)
	}

	return components.RenderPanel(m.Viewport.View(), width, height)
}

// renderCommitMessage renders the commit message panel
func (m Model) renderCommitMessage(width, height int) string {
	content := m.Content
	if m.Err != nil {
		content = components.RenderError(m.Err) + "\n\n" + content
	}

	return components.RenderPanel(content, width, height)
}

// viewEdit renders the edit state
func (m Model) viewEdit() string {
	header := components.RenderHeader("edit commit message", m.Width)
	textarea := components.TextareaStyle.Render(m.Textarea.View())
	hints := renderEditHints(m.Mode)
	stats := m.renderEditorStats()

	return lipgloss.JoinVertical(lipgloss.Left, header, textarea, hints, stats)
}

// renderEditorStats renders the editor statistics (mode, lines, chars)
func (m Model) renderEditorStats() string {
	lines := strings.Count(m.Textarea.Value(), "\n") + 1
	chars := len(m.Textarea.Value())
	stats := fmt.Sprintf("%s  lines: %d  chars: %d", renderModeBadge(m.Mode), lines, chars)
	return components.LabelMedium.Render(stats)
}

// renderEditHints renders keyboard hints for edit mode
func renderEditHints(mode Mode) string {
	if mode == EditModeNormal {
		hints := []string{
			components.KeyHint.Render("[i]") + " insert",
			components.KeyHint.Render("[o]") + " new line",
			components.KeyHint.Render("[d]") + " delete line",
			components.KeyHint.Render("[c]") + " delete word",
			components.KeyHint.Render("[x]") + " delete char",
			components.KeyHint.Render("[s]") + " save",
			components.KeyHint.Render("[←↓↑→]") + " move",
			components.KeyHint.Render("[q]") + " quit",
		}

		return components.RenderHints(hints)
	}

	hints := []string{
		components.KeyHint.Render("[esc]") + " normal",
		components.KeyHint.Render("[ctrl+s]") + " save",
	}

	return components.RenderHints(hints)
}

// renderModeBadge renders the vim mode badge
func renderModeBadge(mode Mode) string {
	if mode == EditModeInsert {
		return components.ModeBadgeInsertStyle.Render(" INSERT ")
	}
	return components.ModeBadgeNormalStyle.Render(" NORMAL ")
}

// viewLogs renders the logs view state
func (m Model) viewLogs() string {
	header := components.RenderHeader("logs", m.Width)

	var content string
	if m.LogBuffer == nil {
		content = components.LabelMedium.Render("Log buffering not enabled (set format: ui in cmt.yaml)")
	} else {
		content = m.Viewport.View()
	}

	hints := components.RenderLogsHints()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, hints)
}
