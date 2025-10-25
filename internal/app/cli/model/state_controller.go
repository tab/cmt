package model

import tea "github.com/charmbracelet/bubbletea"

type stateController interface {
	Handle(m *Model, msg tea.Msg) (tea.Cmd, bool)
}

func defaultControllers() map[State]stateController {
	return map[State]stateController{
		StateFetch:      fetchStateController{},
		StateViewCommit: reviewStateController{},
		StateEditCommit: editStateController{},
		StateViewLogs:   logsStateController{},
	}
}
