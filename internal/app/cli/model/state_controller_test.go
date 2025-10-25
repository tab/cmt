package model

import (
	"context"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"cmt/internal/app/cli/components"
	"cmt/internal/app/cli/workflow"
	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_ModelUpdate(t *testing.T) {
	cfg := config.DefaultConfig()

	cases := []struct {
		name       string
		model      Model
		msg        tea.Msg
		assertFunc func(t *testing.T, next Model, cmd tea.Cmd)
	}{
		{
			name: "window_size_updates_layout",
			model: Model{
				Ctx:      context.Background(),
				Cfg:      cfg,
				Workflow: workflowServiceStub{},
				State:    StateViewCommit,
				Textarea: textarea.New(),
				Viewport: viewport.New(0, 0),

				UserFlow: FlowCommit,
			},
			msg: tea.WindowSizeMsg{Width: 100, Height: 40},
			assertFunc: func(t *testing.T, next Model, cmd tea.Cmd) {
				assert.Nil(t, cmd)
				assert.Equal(t, 100, next.Width)
				assert.Equal(t, 40, next.Height)
				assert.Equal(t, 24, next.Viewport.Width)
				assert.Equal(t, 30, next.Viewport.Height)
			},
		},
		{
			name: "ctrl_c_requests_quit",
			model: Model{
				Ctx:      context.Background(),
				Cfg:      cfg,
				Workflow: workflowServiceStub{},
				State:    StateViewCommit,
				Textarea: textarea.New(),
				Viewport: viewport.New(0, 0),
			},
			msg: tea.KeyMsg{Type: tea.KeyCtrlC},
			assertFunc: func(t *testing.T, next Model, cmd tea.Cmd) {
				assertQuitCmd(t, cmd)
				assert.Equal(t, StateExit, next.State)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, cmd := tc.model.Update(tc.msg)
			next, ok := result.(Model)
			assert.True(t, ok)
			tc.assertFunc(t, next, cmd)
		})
	}
}

func Test_FetchStateController(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(m *Model)
		msg        tea.Msg
		assertFunc func(t *testing.T, m *Model, cmd tea.Cmd, handled bool)
	}{
		{
			name: "propagates_error_and_quits",
			setup: func(m *Model) {
				m.State = StateFetch
			},
			msg: Result{Err: errors.ErrNoGitChanges},
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assertQuitCmd(t, cmd)
				assert.Equal(t, StateExit, m.State)
				assert.ErrorIs(t, m.Err, errors.ErrNoGitChanges)
			},
		},
		{
			name: "transitions_to_commit_view",
			setup: func(m *Model) {
				m.State = StateFetch
				m.UserFlow = FlowCommit
				m.Viewport = viewport.New(20, 10)
			},
			msg: Result{
				Content:  "feat: add feature",
				FileTree: components.ParseGitStatus("A\tfile.go"),
			},
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.Nil(t, cmd)
				assert.Equal(t, StateViewCommit, m.State)
				assert.Equal(t, "feat: add feature", m.Content)
				assert.NotNil(t, m.FileTree)
				assert.NotEmpty(t, m.Viewport.View())
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := minimalModel()
			tc.setup(&model)
			cmd, handled := fetchStateController{}.Handle(&model, tc.msg)
			tc.assertFunc(t, &model, cmd, handled)
		})
	}
}

func Test_ReviewStateController(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(m *Model)
		msg        tea.Msg
		assertFunc func(t *testing.T, m *Model, cmd tea.Cmd, handled bool)
	}{
		{
			name: "accept_triggers_quit",
			setup: func(m *Model) {
				m.State = StateViewCommit
				m.Content = "message"
			},
			msg: runeKey('a'),
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assertQuitCmd(t, cmd)
				assert.Equal(t, ActionAccept, m.UserAction)
				assert.Equal(t, StateExit, m.State)
			},
		},
		{
			name: "edit_enters_editor",
			setup: func(m *Model) {
				m.State = StateViewCommit
				m.Content = "message"
				m.Textarea = textarea.New()
			},
			msg: runeKey('e'),
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.Nil(t, cmd)
				assert.Equal(t, StateEditCommit, m.State)
				assert.Equal(t, ActionEdit, m.UserAction)
				assert.Equal(t, EditModeNormal, m.Mode)
				assert.Equal(t, "message", m.Textarea.Value())
			},
		},
		{
			name: "refresh_triggers_fetch",
			setup: func(m *Model) {
				m.State = StateViewCommit
				m.UserFlow = FlowCommit
				m.Spinner = spinner.New()
				m.Content = "old"
			},
			msg: runeKey('r'),
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.NotNil(t, cmd)
				assert.Equal(t, StateFetch, m.State)
				assert.Empty(t, m.Content)
				assert.Equal(t, ActionRefresh, m.UserAction)
			},
		},
		{
			name: "logs_switches_to_log_state",
			setup: func(m *Model) {
				m.State = StateViewCommit
				m.LogBuffer = logger.NewLogBuffer()
				m.Viewport = viewport.New(20, 10)
			},
			msg: runeKey('l'),
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.Nil(t, cmd)
				assert.Equal(t, StateViewLogs, m.State)
				assert.Equal(t, StateViewCommit, m.PreviousState)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := minimalModel()
			tc.setup(&model)
			cmd, handled := reviewStateController{}.Handle(&model, tc.msg)
			tc.assertFunc(t, &model, cmd, handled)
		})
	}
}

func Test_EditStateController(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(m *Model)
		msg        tea.Msg
		assertFunc func(t *testing.T, m *Model, cmd tea.Cmd, handled bool)
	}{
		{
			name: "normal_mode_switches_to_insert",
			setup: func(m *Model) {
				m.State = StateEditCommit
				m.Textarea = textarea.New()
			},
			msg: runeKey('i'),
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.Nil(t, cmd)
				assert.Equal(t, EditModeInsert, m.Mode)
			},
		},
		{
			name: "insert_mode_ctrl_s_saves_and_returns",
			setup: func(m *Model) {
				m.State = StateEditCommit
				m.Mode = EditModeInsert
				m.Textarea = textarea.New()
				m.Textarea.SetValue("final message")
			},
			msg: tea.KeyMsg{Type: tea.KeyCtrlS},
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.Nil(t, cmd)
				assert.Equal(t, StateViewCommit, m.State)
				assert.Equal(t, "final message", m.Content)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := minimalModel()
			tc.setup(&model)
			cmd, handled := editStateController{}.Handle(&model, tc.msg)
			tc.assertFunc(t, &model, cmd, handled)
		})
	}
}

func Test_LogsStateController(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(m *Model)
		msg        tea.Msg
		assertFunc func(t *testing.T, m *Model, cmd tea.Cmd, handled bool)
	}{
		{
			name: "return_to_previous_state",
			setup: func(m *Model) {
				m.State = StateViewLogs
				m.PreviousState = StateViewCommit
				m.FileTree = components.ParseGitStatus("A\tmain.go")
				m.Viewport = viewport.New(20, 5)
			},
			msg: runeKey('q'),
			assertFunc: func(t *testing.T, m *Model, cmd tea.Cmd, handled bool) {
				assert.True(t, handled)
				assert.Nil(t, cmd)
				assert.Equal(t, StateViewCommit, m.State)
				assert.NotEmpty(t, m.Viewport.View())
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := minimalModel()
			tc.setup(&model)
			cmd, handled := logsStateController{}.Handle(&model, tc.msg)
			tc.assertFunc(t, &model, cmd, handled)
		})
	}
}

type workflowServiceStub struct{}

func (workflowServiceStub) GenerateCommit(ctx context.Context, prefix string) (workflow.CommitResult, error) {
	return workflow.CommitResult{
		Message: "stub message",
	}, nil
}

func (workflowServiceStub) GenerateChangelog(ctx context.Context, between string) (workflow.ChangelogResult, error) {
	return workflow.ChangelogResult{
		Content: "# CHANGELOG",
	}, nil
}

func minimalModel() Model {
	return Model{
		Ctx:      context.Background(),
		Cfg:      config.DefaultConfig(),
		Workflow: workflowServiceStub{},
		Textarea: textarea.New(),
		Viewport: viewport.New(0, 0),

		Log: logger.NewLogger(config.DefaultConfig()),
	}
}

func runeKey(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func assertQuitCmd(t *testing.T, cmd tea.Cmd) {
	t.Helper()
	if !assert.NotNil(t, cmd) {
		return
	}
	msg := cmd()
	_, isQuit := msg.(tea.QuitMsg)
	assert.True(t, isQuit)
}
