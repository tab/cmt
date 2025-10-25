package model

import (
	tea "github.com/charmbracelet/bubbletea"

	"cmt/internal/app/cli/components"
)

type logsStateController struct{}

func (logsStateController) Handle(m *Model, msg tea.Msg) (tea.Cmd, bool) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		switch typed.String() {
		case "q", "l":
			m.State = m.PreviousState
			if m.PreviousState == StateViewCommit && m.FileTree != nil {
				treeContent := components.RenderTree(m.FileTree)
				m.Viewport.SetContent(treeContent)
			}
			return nil, true
		default:
			var cmd tea.Cmd
			m.Viewport, cmd = m.Viewport.Update(typed)
			return cmd, true
		}
	default:
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return cmd, true
	}
}
