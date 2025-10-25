package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"cmt/internal/app/errors"
)

type reviewStateController struct{}

func (reviewStateController) Handle(m *Model, msg tea.Msg) (tea.Cmd, bool) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		switch typed.String() {
		case "a":
			if strings.TrimSpace(m.Content) == "" {
				m.Err = errors.ErrCommitMessageEmpty
				return nil, true
			}
			m.UserAction = ActionAccept
			m.State = StateExit
			return tea.Quit, true
		case "e":
			m.UserAction = ActionEdit
			m.State = StateEditCommit
			m.Mode = EditModeNormal
			m.Textarea.SetValue(m.Content)
			m.Textarea.Focus()
			return nil, true
		case "r":
			if m.UserFlow != FlowCommit {
				return nil, true
			}
			m.UserAction = ActionRefresh
			m.State = StateFetch
			m.Content = ""
			m.Err = nil
			m.FetchStarted = false
			m.Log.Info().Msg("Refreshing commit message")
			return tea.Batch(m.Spinner.Tick, fetchCommit(*m)), true
		case "l":
			if m.LogBuffer != nil {
				m.PreviousState = m.State
				m.State = StateViewLogs
				m.updateLogsViewport()
			}
			return nil, true
		case "q":
			m.UserAction = ActionCancel
			m.State = StateExit
			return tea.Quit, true
		default:
			if m.FileTree != nil {
				var cmd tea.Cmd
				m.Viewport, cmd = m.Viewport.Update(typed)
				return cmd, true
			}
			return nil, true
		}
	default:
		if m.FileTree != nil {
			var cmd tea.Cmd
			m.Viewport, cmd = m.Viewport.Update(msg)
			return cmd, true
		}
		return nil, false
	}
}
