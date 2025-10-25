package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type editStateController struct{}

func (editStateController) Handle(m *Model, msg tea.Msg) (tea.Cmd, bool) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		if m.Mode == EditModeNormal {
			return editStateHandleNormal(m, typed), true
		}
		return editStateHandleInsert(m, typed), true
	default:
		var cmd tea.Cmd
		m.Textarea, cmd = m.Textarea.Update(msg)
		return cmd, true
	}
}

func editStateHandleNormal(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	case "i":
		m.Mode = EditModeInsert
	case "o":
		m.Mode = EditModeInsert
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyEnd})
		m.Textarea.InsertString("\n")
	case "h", "left":
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyLeft})
	case "j", "down":
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyDown})
	case "k", "up":
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyUp})
	case "l", "right":
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyRight})
	case "c":
		m.Textarea, _ = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyCtrlLeft})

		for i := 0; i < maxDeleteIterations; i++ {
			before := m.Textarea.Value()
			if len(before) == 0 {
				break
			}

			m.Textarea, _ = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyDelete})
			after := m.Textarea.Value()

			if before == after {
				break
			}

			deletedChar := byte(0)
			for j := 0; j < len(before) && j < len(after); j++ {
				if before[j] != after[j] {
					deletedChar = before[j]
					break
				}
			}
			if deletedChar == 0 && len(after) < len(before) {
				deletedChar = before[len(after)]
			}

			if deletedChar == ' ' || deletedChar == '\n' || deletedChar == '\t' {
				break
			}
		}
	case "x":
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyDelete})
	case "d":
		m.Textarea, _ = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyHome})
		m.Textarea, _ = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyCtrlK})
		m.Textarea, cmd = m.Textarea.Update(tea.KeyMsg{Type: tea.KeyDelete})
	case "s":
		m.Content = m.Textarea.Value()
		if strings.TrimSpace(m.Content) != "" {
			m.Err = nil
		}
		m.State = StateViewCommit
		m.Textarea.Blur()
	case "q":
		m.State = StateViewCommit
		m.Textarea.Blur()
	}

	return cmd
}

func editStateHandleInsert(m *Model, msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEsc:
		m.Mode = EditModeNormal
		return nil
	case tea.KeyCtrlS:
		m.Content = m.Textarea.Value()
		if strings.TrimSpace(m.Content) != "" {
			m.Err = nil
		}
		m.State = StateViewCommit
		m.Textarea.Blur()
		return nil
	default:
		var cmd tea.Cmd
		m.Textarea, cmd = m.Textarea.Update(msg)
		return cmd
	}
}
