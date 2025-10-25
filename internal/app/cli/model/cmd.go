package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Start initiates the commit workflow
func Start(m Model) tea.Cmd {
	return fetchCommit(m)
}

// fetchCommit fetches git diff and generates commit message
func fetchCommit(m Model) tea.Cmd {
	return func() tea.Msg {
		result, err := m.Workflow.GenerateCommit(m.Ctx, m.Prefix)
		if err != nil {
			return Result{Err: err}
		}

		return Result{Content: result.Message, FileTree: result.FileTree}
	}
}
