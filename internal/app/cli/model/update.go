package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const maxDeleteIterations = 100

// Update delegates Bubble Tea messages to state-specific controllers.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m = m.ensureControllers()

	if typed, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width = typed.Width
		m.Height = typed.Height

		leftWidth := int(float64(typed.Width) * 0.3)
		m.Viewport.Width = leftWidth - 6
		m.Viewport.Height = typed.Height - 10

		m.Textarea.SetWidth(typed.Width - 4)
		m.Textarea.SetHeight(typed.Height - 12)

		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyCtrlC {
		m.State = StateExit
		return m, tea.Quit
	}

	if controller, ok := m.controllers[m.State]; ok {
		if cmd, handled := controller.Handle(&m, msg); handled {
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) ensureControllers() Model {
	if m.controllers == nil {
		m.controllers = defaultControllers()
	}
	return m
}

// updateLogsViewport updates the logs viewport with current log entries.
func (m *Model) updateLogsViewport() {
	if m.LogBuffer == nil {
		return
	}

	m.Viewport.Width = m.Width - 4
	m.Viewport.Height = m.Height - 6

	entries := m.LogBuffer.GetEntries()
	lines := make([]string, 0, len(entries))
	for i := range entries {
		lines = append(lines, FormatLogEntry(&entries[i]))
	}

	m.Viewport.SetContent(strings.Join(lines, "\n"))
}
