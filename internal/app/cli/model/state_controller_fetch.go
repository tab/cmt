package model

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"cmt/internal/app/cli/components"
)

type fetchStateController struct{}

func (fetchStateController) Handle(m *Model, msg tea.Msg) (tea.Cmd, bool) {
	switch typed := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(typed)
		if !m.FetchStarted && m.Content == "" && m.Err == nil {
			m.FetchStarted = true
			return tea.Batch(cmd, Start(*m)), true
		}
		return cmd, true
	case Result:
		if typed.Err != nil {
			m.Err = typed.Err
			m.State = StateExit
			m.FetchStarted = false
			return tea.Quit, true
		}

		m.Content = typed.Content
		if typed.FileTree != nil {
			m.FileTree = typed.FileTree
		}

		m.State = StateViewCommit
		if typed.FileTree != nil {
			treeContent := components.RenderTree(typed.FileTree)
			m.Viewport.SetContent(treeContent)
		}

		return nil, true
	case tea.KeyMsg:
		// Ignore key input while fetching.
		return nil, true
	default:
		return nil, false
	}
}
