package commit

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

func Test_NewModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name          string
		commitMessage string
		prefix        string
		expectedMode  WorkflowMode
	}{
		{
			name:          "Success with commit message",
			commitMessage: "Initial commit",
			prefix:        "feat:",
			expectedMode:  Viewing,
		},
		{
			name:          "Success without commit message",
			commitMessage: "",
			prefix:        "",
			expectedMode:  Fetching,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				Files:         "A\tfile.txt",
				CommitMessage: tt.commitMessage,
				Prefix:        tt.prefix,
				Diff:          "some diff",
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)

			assert.Equal(t, tt.expectedMode, m.stateMachine.WorkflowMode())
			assert.Equal(t, tt.commitMessage, m.state.CommitMessage)
			assert.Equal(t, tt.prefix, m.state.Prefix)
			assert.Equal(t, MessageFocus, m.focusPane)
		})
	}
}

func Test_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name          string
		commitMessage string
		expectCmd     bool
	}{
		{
			name:          "Success with viewing mode command",
			commitMessage: "test",
			expectCmd:     true,
		},
		{
			name:          "Success with fetching mode command",
			commitMessage: "",
			expectCmd:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: tt.commitMessage,
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			cmd := m.Init()

			if tt.expectCmd {
				assert.NotNil(t, cmd)
			}
		})
	}
}

func Test_FetchInitialData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	ctx := context.Background()

	tests := []struct {
		name        string
		before      func(*git.MockClient, *gpt.MockClient)
		expectError bool
		checkFn     func(*testing.T, tea.Msg)
	}{
		{
			name: "Success",
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().Status(ctx).Return("A\tfile.txt", nil)
				mockGit.EXPECT().Diff(ctx).Return("diff content", nil)
				mockGPT.EXPECT().FetchCommitMessage(ctx, "diff content").Return("Generated message", nil)
			},
			expectError: false,
			checkFn: func(t *testing.T, msg tea.Msg) {
				successMsg, ok := msg.(FetchSuccessMsg)
				assert.True(t, ok)
				assert.Equal(t, "A\tfile.txt", successMsg.Status)
				assert.Equal(t, "diff content", successMsg.Diff)
				assert.Equal(t, "Generated message", successMsg.Message)
			},
		},
		{
			name: "Failure when git status fails",
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().Status(ctx).Return("", errors.New("git status failed"))
			},
			expectError: true,
			checkFn: func(t *testing.T, msg tea.Msg) {
				errorMsg, ok := msg.(FetchErrorMsg)
				assert.True(t, ok)
				assert.Error(t, errorMsg.Err)
			},
		},
		{
			name: "Failure when git diff fails",
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().Status(ctx).Return("A\tfile.txt", nil)
				mockGit.EXPECT().Diff(ctx).Return("", errors.New("git diff failed"))
			},
			expectError: true,
			checkFn: func(t *testing.T, msg tea.Msg) {
				errorMsg, ok := msg.(FetchErrorMsg)
				assert.True(t, ok)
				assert.Error(t, errorMsg.Err)
			},
		},
		{
			name: "Failure when g p t fails",
			before: func(mockGit *git.MockClient, mockGPT *gpt.MockClient) {
				mockGit.EXPECT().Status(ctx).Return("A\tfile.txt", nil)
				mockGit.EXPECT().Diff(ctx).Return("diff content", nil)
				mockGPT.EXPECT().FetchCommitMessage(ctx, "diff content").Return("", errors.New("gpt fetch failed"))
			},
			expectError: true,
			checkFn: func(t *testing.T, msg tea.Msg) {
				errorMsg, ok := msg.(FetchErrorMsg)
				assert.True(t, ok)
				assert.Error(t, errorMsg.Err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGit := git.NewMockClient(ctrl)
			mockGPT := gpt.NewMockClient(ctrl)
			mockSpinner := spinner.NewMockModel(ctrl)

			mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

			tt.before(mockGit, mockGPT)

			input := Input{
				GitClient: mockGit,
				GPTClient: mockGPT,
				Logger:    mockLogger,
				Ctx:       ctx,
				Spinner:   func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			cmd := m.fetchInitialData()
			msg := cmd()

			tt.checkFn(t, msg)
		})
	}
}

func Test_Update_WindowSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		Files:         "A\tfile.txt",
		CommitMessage: "test",
		GitClient:     mockGit,
		GPTClient:     mockGPT,
		Logger:        mockLogger,
		Ctx:           context.Background(),
		Spinner:       func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updated, _ := m.Update(msg)
	updatedModel := updated.(Model)

	assert.Equal(t, 120, updatedModel.width)
	assert.Equal(t, 40, updatedModel.height)
	assert.True(t, updatedModel.ready)
}

func Test_Update_FetchSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		GitClient: mockGit,
		GPTClient: mockGPT,
		Logger:    mockLogger,
		Ctx:       context.Background(),
		Spinner:   func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	msg := FetchSuccessMsg{
		Status:  "A\tfile.txt",
		Diff:    "diff content",
		Message: "Generated message",
	}

	updated, _ := m.Update(msg)
	updatedModel := updated.(Model)

	assert.Equal(t, 1, len(updatedModel.state.Files))
	assert.Equal(t, "diff content", updatedModel.state.Diff)
	assert.Equal(t, "Generated message", updatedModel.state.CommitMessage)
	assert.Equal(t, Viewing, updatedModel.stateMachine.WorkflowMode())
}

func Test_Update_FetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		GitClient: mockGit,
		GPTClient: mockGPT,
		Logger:    mockLogger,
		Ctx:       context.Background(),
		Spinner:   func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	expectedErr := errors.New("fetch failed")
	msg := FetchErrorMsg{Err: expectedErr}

	updated, cmd := m.Update(msg)
	updatedModel := updated.(Model)

	assert.Equal(t, expectedErr, updatedModel.state.Error)
	assert.Equal(t, Viewing, updatedModel.stateMachine.WorkflowMode())
	assert.NotNil(t, cmd)
}

func Test_Update_RegenerateMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()
	mockSpinner.EXPECT().Update(gomock.Any()).Return(mockSpinner, nil).AnyTimes()

	tests := []struct {
		name            string
		oldMessage      string
		regenerateMsg   RegenerateMsg
		expectedMessage string
		expectedMode    WorkflowMode
	}{
		{
			name:       "Success when message updates",
			oldMessage: "old message",
			regenerateMsg: RegenerateMsg{
				Message: "new message",
				Err:     nil,
			},
			expectedMessage: "new message",
			expectedMode:    Viewing,
		},
		{
			name:       "Failure when message update fails",
			oldMessage: "old message",
			regenerateMsg: RegenerateMsg{
				Message: "",
				Err:     errors.New("regenerate failed"),
			},
			expectedMessage: "old message",
			expectedMode:    Viewing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: tt.oldMessage,
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.width = 120
			m.height = 40
			m.ready = true

			updated, _ := m.Update(tt.regenerateMsg)
			updatedModel := updated.(Model)

			assert.Equal(t, tt.expectedMessage, updatedModel.state.CommitMessage)
			assert.Equal(t, tt.expectedMode, updatedModel.stateMachine.WorkflowMode())
		})
	}
}

func Test_HandleNormalMode_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		CommitMessage: "test commit",
		GitClient:     mockGit,
		GPTClient:     mockGPT,
		Logger:        mockLogger,
		Ctx:           ctx,
		Spinner:       func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updated, cmd := m.handleNormalMode(keyMsg)
	updatedModel := updated.(Model)

	assert.True(t, updatedModel.state.Accepted)
	assert.NotNil(t, cmd)
}

func Test_HandleNormalMode_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		CommitMessage: "test commit",
		GitClient:     mockGit,
		GPTClient:     mockGPT,
		Logger:        mockLogger,
		Ctx:           context.Background(),
		Spinner:       func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	updated, _ := m.handleNormalMode(keyMsg)
	updatedModel := updated.(Model)

	assert.Equal(t, Editing, updatedModel.stateMachine.WorkflowMode())
	assert.Equal(t, "test commit", updatedModel.textarea.Value())
}

func Test_HandleNormalMode_Edit_WithPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name          string
		commitMessage string
		prefix        string
		expectedValue string
	}{
		{
			name:          "Edit mode with prefix",
			commitMessage: "chore(ui): Add login feature",
			prefix:        "JIRA-123",
			expectedValue: "JIRA-123 chore(ui): Add login feature",
		},
		{
			name:          "Edit mode without prefix",
			commitMessage: "chore(ui): Add login feature",
			prefix:        "",
			expectedValue: "chore(ui): Add login feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: tt.commitMessage,
				Prefix:        tt.prefix,
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			updated, _ := m.handleNormalMode(keyMsg)
			updatedModel := updated.(Model)

			assert.Equal(t, Editing, updatedModel.stateMachine.WorkflowMode())
			assert.Equal(t, tt.expectedValue, updatedModel.textarea.Value())
		})
	}
}

func Test_HandleNormalMode_Regenerate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()
	mockGPT.EXPECT().FetchCommitMessage(ctx, "test diff").Return("regenerated message", nil)

	input := Input{
		CommitMessage: "old message",
		Diff:          "test diff",
		GitClient:     mockGit,
		GPTClient:     mockGPT,
		Logger:        mockLogger,
		Ctx:           ctx,
		Spinner:       func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updated, cmd := m.handleNormalMode(keyMsg)
	updatedModel := updated.(Model)

	assert.Equal(t, Regenerating, updatedModel.stateMachine.WorkflowMode())
	assert.NotNil(t, cmd)

	// Execute the returned command
	msg := cmd()
	regenMsg, ok := msg.(RegenerateMsg)
	assert.True(t, ok)
	assert.Equal(t, "regenerated message", regenMsg.Message)
}

func Test_HandleNormalMode_ToggleLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name         string
		initialPane  ViewPane
		expectedPane ViewPane
	}{
		{
			name:         "Success when toggling message to logs",
			initialPane:  MessagePane,
			expectedPane: AppLogsPane,
		},
		{
			name:         "Success when toggling logs to message",
			initialPane:  AppLogsPane,
			expectedPane: MessagePane,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := logger.NewMockLogger(ctrl)

			// When toggling to AppLogsPane, renderAppLogs will call GetBuffer
			if tt.expectedPane == AppLogsPane {
				buffer := logger.NewLogBuffer(100)
				mockLogger.EXPECT().GetBuffer().Return(buffer)
			}

			input := Input{
				GitClient: mockGit,
				GPTClient: mockGPT,
				Logger:    mockLogger,
				Ctx:       context.Background(),
				Spinner:   func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.width = 120
			m.height = 40
			m.ready = true
			m.stateMachine.EnterViewing(tt.initialPane)

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			updated, _ := m.handleNormalMode(keyMsg)
			updatedModel := updated.(Model)

			assert.Equal(t, tt.expectedPane, updatedModel.stateMachine.ViewPane())
		})
	}
}

func Test_HandleNormalMode_ToggleFocus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		GitClient: mockGit,
		GPTClient: mockGPT,
		Logger:    mockLogger,
		Ctx:       context.Background(),
		Spinner:   func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	m.focusPane = TreeFocus

	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.handleNormalMode(keyMsg)
	updatedModel := updated.(Model)

	assert.Equal(t, MessageFocus, updatedModel.focusPane)

	// Toggle back
	updated2, _ := updatedModel.handleNormalMode(keyMsg)
	updatedModel2 := updated2.(Model)

	assert.Equal(t, TreeFocus, updatedModel2.focusPane)
}

func Test_HandleNormalMode_Quit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		GitClient: mockGit,
		GPTClient: mockGPT,
		Logger:    mockLogger,
		Ctx:       context.Background(),
		Spinner:   func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.handleNormalMode(keyMsg)

	assert.NotNil(t, cmd)
}

func Test_HandleEditMode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name    string
		keyType tea.KeyType
		checkFn func(*testing.T, Model, tea.Cmd)
	}{
		{
			name:    "Success with escape saves",
			keyType: tea.KeyEsc,
			checkFn: func(t *testing.T, m Model, cmd tea.Cmd) {
				assert.Equal(t, "new message", m.state.CommitMessage)
				assert.Equal(t, Viewing, m.stateMachine.WorkflowMode())
			},
		},
		{
			name:    "Success with ctrl c quits",
			keyType: tea.KeyCtrlC,
			checkFn: func(t *testing.T, m Model, cmd tea.Cmd) {
				assert.NotNil(t, cmd)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: "old",
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.width = 120
			m.height = 40
			m.ready = true
			m.stateMachine.EnterEditing()
			m.textarea.SetValue("new message")

			keyMsg := tea.KeyMsg{Type: tt.keyType}
			updated, cmd := m.handleEditMode(keyMsg)
			updatedModel := updated.(Model)

			tt.checkFn(t, updatedModel, cmd)
		})
	}
}

func Test_HandleEditMode_WithPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name           string
		originalPrefix string
		editedText     string
		expectedPrefix string
		expectedMsg    string
	}{
		{
			name:           "Keeps original prefix when present",
			originalPrefix: "JIRA-123",
			editedText:     "JIRA-123 chore(ui): Updated commit message",
			expectedPrefix: "JIRA-123",
			expectedMsg:    "chore(ui): Updated commit message",
		},
		{
			name:           "Handles removed prefix",
			originalPrefix: "JIRA-123",
			editedText:     "chore(ui): Plain commit message",
			expectedPrefix: "",
			expectedMsg:    "chore(ui): Plain commit message",
		},
		{
			name:           "Handles no original prefix",
			originalPrefix: "",
			editedText:     "chore(ui): Plain commit message",
			expectedPrefix: "",
			expectedMsg:    "chore(ui): Plain commit message",
		},
		{
			name:           "Handles changed prefix",
			originalPrefix: "OLD-123",
			editedText:     "NEW-456 chore(ui): Changed message",
			expectedPrefix: "",
			expectedMsg:    "NEW-456 chore(ui): Changed message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: "old message",
				Prefix:        tt.originalPrefix,
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.width = 120
			m.height = 40
			m.ready = true
			m.stateMachine.EnterEditing()
			m.textarea.SetValue(tt.editedText)

			keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
			updated, _ := m.handleEditMode(keyMsg)
			updatedModel := updated.(Model)

			assert.Equal(t, tt.expectedPrefix, updatedModel.state.Prefix)
			assert.Equal(t, tt.expectedMsg, updatedModel.state.CommitMessage)
			assert.Equal(t, Viewing, updatedModel.stateMachine.WorkflowMode())
		})
	}
}

func Test_GetOutput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name          string
		commitMessage string
		prefix        string
		accepted      bool
		result        string
	}{
		{
			name:          "Success without prefix",
			commitMessage: "test commit",
			prefix:        "",
			accepted:      true,
			result:        "test commit",
		},
		{
			name:          "Success with prefix",
			commitMessage: "test commit",
			prefix:        "feat:",
			accepted:      true,
			result:        "feat: test commit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: tt.commitMessage,
				Prefix:        tt.prefix,
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.state.Accepted = tt.accepted

			output := m.GetOutput()

			assert.Equal(t, tt.result, output.Result)
			assert.Equal(t, tt.accepted, output.Accepted)
		})
	}
}

func Test_GetDisplayMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name          string
		commitMessage string
		prefix        string
		ready         bool
		checkFn       func(*testing.T, string)
	}{
		{
			name:          "Success without prefix",
			commitMessage: "test commit",
			prefix:        "",
			ready:         false,
			checkFn: func(t *testing.T, msg string) {
				assert.Equal(t, "test commit", msg)
			},
		},
		{
			name:          "Success with prefix",
			commitMessage: "test commit",
			prefix:        "feat:",
			ready:         false,
			checkFn: func(t *testing.T, msg string) {
				assert.Contains(t, msg, "feat: test commit")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: tt.commitMessage,
				Prefix:        tt.prefix,
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.ready = tt.ready

			msg := m.getDisplayMessage()
			tt.checkFn(t, msg)
		})
	}
}

func Test_RenderAppLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name     string
		setupLog func(*gomock.Controller) logger.Logger
		expected string
	}{
		{
			name: "Success with nil logger",
			setupLog: func(ctrl *gomock.Controller) logger.Logger {
				return nil
			},
			expected: "No logger available",
		},
		{
			name: "Success with nil buffer",
			setupLog: func(ctrl *gomock.Controller) logger.Logger {
				mockLogger := logger.NewMockLogger(ctrl)
				mockLogger.EXPECT().GetBuffer().Return(nil)
				return mockLogger
			},
			expected: "Log buffer not available (console mode)",
		},
		{
			name: "Success with empty entries",
			setupLog: func(ctrl *gomock.Controller) logger.Logger {
				mockLogger := logger.NewMockLogger(ctrl)
				buffer := logger.NewLogBuffer(100)
				mockLogger.EXPECT().GetBuffer().Return(buffer)
				return mockLogger
			},
			expected: "No logs yet",
		},
		{
			name: "Success with formatted entries",
			setupLog: func(ctrl *gomock.Controller) logger.Logger {
				mockLogger := logger.NewMockLogger(ctrl)
				buffer := logger.NewLogBuffer(100)
				buffer.AddFormatted("INFO", "test message", "formatted line 1")
				buffer.AddFormatted("ERROR", "error message", "formatted line 2")
				mockLogger.EXPECT().GetBuffer().Return(buffer)
				return mockLogger
			},
			expected: "formatted line 1",
		},
		{
			name: "Success with unformatted entries",
			setupLog: func(ctrl *gomock.Controller) logger.Logger {
				mockLogger := logger.NewMockLogger(ctrl)
				buffer := logger.NewLogBuffer(100)
				buffer.Add("INFO", "plain message")
				mockLogger.EXPECT().GetBuffer().Return(buffer)
				return mockLogger
			},
			expected: "plain message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				GitClient: mockGit,
				GPTClient: mockGPT,
				Logger:    tt.setupLog(ctrl),
				Ctx:       context.Background(),
				Spinner:   func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			logs := m.renderAppLogs()

			assert.Contains(t, logs, tt.expected)
		})
	}
}

func Test_ParsePrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name           string
		text           string
		originalPrefix string
		expectedPrefix string
		expectedMsg    string
	}{
		{
			name:           "Original prefix with conventional commit",
			text:           "JIRA-123 chore(ui): commit message",
			originalPrefix: "JIRA-123",
			expectedPrefix: "JIRA-123",
			expectedMsg:    "chore(ui): commit message",
		},
		{
			name:           "Task prefix with conventional commit",
			text:           "TASK-12345 feat(auth): add login",
			originalPrefix: "TASK-12345",
			expectedPrefix: "TASK-12345",
			expectedMsg:    "feat(auth): add login",
		},
		{
			name:           "User removed prefix",
			text:           "chore(ui): commit message",
			originalPrefix: "JIRA-123",
			expectedPrefix: "",
			expectedMsg:    "chore(ui): commit message",
		},
		{
			name:           "User changed prefix",
			text:           "NEW-123 chore(ui): message",
			originalPrefix: "OLD-123",
			expectedPrefix: "",
			expectedMsg:    "NEW-123 chore(ui): message",
		},
		{
			name:           "No original prefix set",
			text:           "chore(ui): some message",
			originalPrefix: "",
			expectedPrefix: "",
			expectedMsg:    "chore(ui): some message",
		},
		{
			name:           "Empty text",
			text:           "",
			originalPrefix: "",
			expectedPrefix: "",
			expectedMsg:    "",
		},
		{
			name:           "Text with whitespace only",
			text:           "   ",
			originalPrefix: "",
			expectedPrefix: "",
			expectedMsg:    "",
		},
		{
			name:           "Prefix without space",
			text:           "JIRA-123chore",
			originalPrefix: "JIRA-123",
			expectedPrefix: "",
			expectedMsg:    "JIRA-123chore",
		},
		{
			name:           "Prefix with extra spaces",
			text:           "JIRA-123  chore(ui): message",
			originalPrefix: "JIRA-123",
			expectedPrefix: "JIRA-123",
			expectedMsg:    "chore(ui): message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				Prefix:    tt.originalPrefix,
				GitClient: mockGit,
				GPTClient: mockGPT,
				Logger:    mockLogger,
				Ctx:       context.Background(),
				Spinner:   func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			prefix, message := m.parsePrefix(tt.text)

			assert.Equal(t, tt.expectedPrefix, prefix)
			assert.Equal(t, tt.expectedMsg, message)
		})
	}
}
