package model

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"cmt/internal/config"
	"cmt/internal/config/logger"
)

// New creates a new UI model with the given dependencies
func New(ctx context.Context, cfg *config.Config, workflow workflowService, log logger.Logger, args []string) Model {
	prefix := parseArgs(args)

	logEvent := log.Info().Str("flow", "commit")
	if prefix != "" {
		logEvent = logEvent.Str("prefix", prefix)
	}
	logEvent.Msg("Initializing workflow")

	s := spinner.New()
	s.Spinner = spinner.Dot

	ta := textarea.New()
	ta.Placeholder = "enter commit message..."
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = true

	vp := viewport.New(0, 0)

	var logBuffer *logger.LogBuffer
	if appLogger, ok := log.(*logger.AppLogger); ok {
		logBuffer = appLogger.GetBuffer()
	}

	return Model{
		Ctx:           ctx,
		Cfg:           cfg,
		Workflow:      workflow,
		State:         StateFetch,
		Spinner:       s,
		Textarea:      ta,
		Viewport:      vp,
		LogBuffer:     logBuffer,
		PreviousState: StateInit,
		Content:       "",
		FileTree:      nil,
		Err:           nil,
		Width:         0,
		Height:        0,
		UserFlow:      FlowCommit,
		UserAction:    ActionNone,
		Prefix:        prefix,
		Log:           log,
	}
}

// Init initializes the model (Bubble Tea lifecycle)
func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

// parseArgs extracts prefix from args
func parseArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	cmd := args[0]

	if Contains(CmdPrefix, cmd) {
		if len(args) > 1 {
			return args[1]
		}
		return ""
	}

	return args[0]
}
