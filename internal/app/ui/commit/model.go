package commit

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

// FocusPane represents which pane has focus
type FocusPane int

const (
	// TreeFocus indicates the file tree pane has focus
	TreeFocus FocusPane = iota
	// MessageFocus indicates the message pane has focus
	MessageFocus
)

// Model represents the Bubble Tea model for commit UI
type Model struct {
	state        State
	stateMachine *stateMachine
	keys         KeyMap
	help         help.Model
	viewport     viewport.Model
	treeViewport viewport.Model
	textarea     textarea.Model
	spinner      spinner.Model
	gitClient    git.Client
	gptClient    gpt.Client
	logger       logger.Logger
	ctx          context.Context
	width        int
	height       int
	ready        bool
	focusPane    FocusPane
}

// Input contains the initial data for the commit UI
type Input struct {
	Files         string
	CommitMessage string
	Prefix        string
	Diff          string
	GitClient     git.Client
	GPTClient     gpt.Client
	Logger        logger.Logger
	Ctx           context.Context
	Spinner       spinner.Factory
}

// Output contains the result after the commit UI exits
type Output struct {
	Accepted      bool
	CommitMessage string
	CommitOutput  string
	Error         error
}

// NewModel creates a new commit UI model
func NewModel(input Input) Model {
	keys := DefaultKeyMap()
	h := help.New()
	h.ShowAll = false

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle()
	treeVp := viewport.New(0, 0)
	ta := textarea.New()
	ta.Placeholder = "Enter commit message…"
	ta.CharLimit = 0

	s := input.Spinner()

	files := BuildFileTree(input.Files)

	initialMode := Viewing
	if input.CommitMessage == "" {
		initialMode = Fetching
	}

	return Model{
		state: State{
			Files:         files,
			CommitMessage: input.CommitMessage,
			Prefix:        input.Prefix,
			Diff:          input.Diff,
		},
		stateMachine: newStateMachine(MessagePane, initialMode),
		keys:         keys,
		help:         h,
		viewport:     vp,
		treeViewport: treeVp,
		textarea:     ta,
		spinner:      s,
		gitClient:    input.GitClient,
		gptClient:    input.GPTClient,
		logger:       input.Logger,
		ctx:          input.Ctx,
		ready:        false,
		focusPane:    MessageFocus,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	if m.stateMachine.WorkflowMode() == Fetching {
		return tea.Batch(m.spinner.Tick, m.fetchInitialData())
	}
	return m.spinner.Tick
}

// fetchInitialData fetches git status, diff, and generates initial commit message
func (m Model) fetchInitialData() tea.Cmd {
	return func() tea.Msg {
		status, err := m.gitClient.Status(m.ctx)
		if err != nil {
			return FetchErrorMsg{Err: err}
		}

		diff, err := m.gitClient.Diff(m.ctx)
		if err != nil {
			return FetchErrorMsg{Err: err}
		}

		message, err := m.gptClient.FetchCommitMessage(m.ctx, diff)
		if err != nil {
			return FetchErrorMsg{Err: err}
		}

		return FetchSuccessMsg{
			Status:  status,
			Diff:    diff,
			Message: message,
		}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		treeWidth := m.width / 3
		messageWidth := m.width - treeWidth - 4

		panelHeight := m.height - 10
		viewportHeight := panelHeight - 2

		m.treeViewport.Width = treeWidth - 4
		m.treeViewport.Height = viewportHeight
		m.treeViewport.SetContent(m.renderFileTree())

		switch m.stateMachine.ViewPane() {
		case AppLogsPane:
			m.viewport.Width = m.width - 6
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(m.renderAppLogs())
		default:
			m.viewport.Width = messageWidth - 4
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(m.getDisplayMessage())
		}

		m.textarea.SetWidth(m.width)
		m.textarea.SetHeight(m.height - 5)

		return m, nil

	case tea.KeyMsg:
		if m.stateMachine.WorkflowMode() == Editing {
			return m.handleEditMode(msg)
		}
		return m.handleNormalMode(msg)

	case FetchSuccessMsg:
		files := BuildFileTree(msg.Status)
		m.state.Files = files
		m.state.Diff = msg.Diff
		m.state.CommitMessage = msg.Message
		m.viewport.SetContent(m.getDisplayMessage())
		m.treeViewport.SetContent(m.renderFileTree())
		m.stateMachine.EnterViewing(MessagePane)
		return m, nil

	case FetchErrorMsg:
		m.state.Error = msg.Err
		m.stateMachine.EnterViewing(MessagePane)
		return m, tea.Quit

	case RegenerateMsg:
		if msg.Err != nil {
			m.stateMachine.EnterViewing(m.stateMachine.ViewPane())
		} else {
			m.state.CommitMessage = msg.Message
			if m.stateMachine.ViewPane() == MessagePane {
				m.viewport.SetContent(m.getDisplayMessage())
			}
			m.stateMachine.EnterViewing(m.stateMachine.ViewPane())
		}
		return m, nil

	case CommitSuccessMsg:
		m.state.Accepted = true
		m.state.CommitOutput = msg.Output
		return m, tea.Quit

	case CommitErrorMsg:
		m.state.Error = msg.Err
		m.stateMachine.EnterViewing(m.stateMachine.ViewPane())
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// handleNormalMode processes keys in normal (viewing) mode
func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Accept):
		if m.stateMachine.CanAccept() {
			m.stateMachine.EnterCommitting()
			return m, m.acceptAndCommit()
		}
		return m, nil

	case key.Matches(msg, m.keys.Edit):
		if m.stateMachine.CanEdit() {
			m.stateMachine.EnterEditing()
			m.textarea.SetValue(m.state.CommitMessage)
			m.textarea.Focus()
			return m, textarea.Blink
		}
		return m, nil

	case key.Matches(msg, m.keys.Regenerate):
		if m.stateMachine.CanRegenerate() {
			m.stateMachine.EnterRegenerating()
			return m, m.regenerateMessage()
		}
		return m, nil

	case key.Matches(msg, m.keys.ToggleLogs):
		if m.stateMachine.CanToggleView() {
			if m.stateMachine.ViewPane() == AppLogsPane {
				m.stateMachine.EnterViewing(MessagePane)
				treeWidth := m.width / 3
				messageWidth := m.width - treeWidth - 4
				m.viewport.Width = messageWidth - 4
				m.viewport.SetContent(m.getDisplayMessage())
			} else {
				m.stateMachine.EnterViewing(AppLogsPane)
				m.viewport.Width = m.width - 6
				m.viewport.SetContent(m.renderAppLogs())
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.ToggleFocus):
		if m.focusPane == TreeFocus {
			m.focusPane = MessageFocus
		} else {
			m.focusPane = TreeFocus
		}
		return m, nil

	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	}

	switch msg.String() {
	case "j", "down", "k", "up", "pgdown", "pgup", "home", "end", "g", "G":
		if m.stateMachine.ViewPane() == AppLogsPane {
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		} else {
			if m.focusPane == TreeFocus {
				m.treeViewport, cmd = m.treeViewport.Update(msg)
			} else {
				m.viewport, cmd = m.viewport.Update(msg)
			}
			return m, cmd
		}
	}

	return m, nil
}

// handleEditMode processes keys in edit mode
func (m Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.state.CommitMessage = m.textarea.Value()
		pane := m.stateMachine.ViewPane()
		if pane == AppLogsPane {
			pane = m.stateMachine.LastViewPane()
		}
		m.stateMachine.EnterViewing(pane)
		if m.stateMachine.ViewPane() == MessagePane {
			m.viewport.SetContent(m.getDisplayMessage())
		}
		m.textarea.Blur()
		return m, nil

	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// acceptAndCommit creates a command to commit the current message
func (m Model) acceptAndCommit() tea.Cmd {
	return func() tea.Msg {
		message := m.state.CommitMessage
		if m.state.Prefix != "" {
			message = m.state.Prefix + " " + message
		}

		output, err := m.gitClient.Commit(m.ctx, message)
		if err != nil {
			return CommitErrorMsg{Err: err}
		}

		return CommitSuccessMsg{Output: output}
	}
}

// regenerateMessage creates a command to regenerate the commit message
func (m Model) regenerateMessage() tea.Cmd {
	return func() tea.Msg {
		message, err := m.gptClient.FetchCommitMessage(m.ctx, m.state.Diff)
		return RegenerateMsg{Message: message, Err: err}
	}
}

// GetOutput returns the final output after the program exits
func (m Model) GetOutput() Output {
	message := m.state.CommitMessage
	if m.state.Prefix != "" {
		message = m.state.Prefix + " " + message
	}

	return Output{
		Accepted:      m.state.Accepted,
		CommitMessage: message,
		CommitOutput:  m.state.CommitOutput,
		Error:         m.state.Error,
	}
}

// getDisplayMessage returns the commit message with prefix prepended if set
func (m Model) getDisplayMessage() string {
	message := m.state.CommitMessage
	if m.state.Prefix != "" {
		message = m.state.Prefix + " " + message
	}

	if m.ready && m.viewport.Width > 0 {
		return lipgloss.NewStyle().Width(m.viewport.Width).Render(message)
	}

	return message
}

// renderAppLogs retrieves and formats application logs from the buffer
func (m Model) renderAppLogs() string {
	if m.logger == nil {
		return "No logger available"
	}

	buffer := m.logger.GetBuffer()
	if buffer == nil {
		return "Log buffer not available (console mode)"
	}

	entries := buffer.Entries()
	if len(entries) == 0 {
		return "No logs yet"
	}

	var sb strings.Builder
	for _, entry := range entries {
		if entry.FormattedLine != "" {
			sb.WriteString(entry.FormattedLine)
			sb.WriteString("\n")
		} else {
			sb.WriteString(entry.Message)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
