package commit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ViewPane_String(t *testing.T) {
	tests := []struct {
		name     string
		pane     ViewPane
		expected string
	}{
		{
			name:     "Success with message pane",
			pane:     MessagePane,
			expected: "MessagePane",
		},
		{
			name:     "Success with app logs pane",
			pane:     AppLogsPane,
			expected: "AppLogsPane",
		},
		{
			name:     "Success with unknown pane",
			pane:     ViewPane(999),
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pane.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_WorkflowMode_String(t *testing.T) {
	tests := []struct {
		name     string
		mode     WorkflowMode
		expected string
	}{
		{
			name:     "Success with viewing mode",
			mode:     Viewing,
			expected: "Viewing",
		},
		{
			name:     "Success with editing mode",
			mode:     Editing,
			expected: "Editing",
		},
		{
			name:     "Success with fetching mode",
			mode:     Fetching,
			expected: "Fetching",
		},
		{
			name:     "Success with regenerating mode",
			mode:     Regenerating,
			expected: "Regenerating",
		},
		{
			name:     "Success with unknown mode",
			mode:     WorkflowMode(999),
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mode.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_NewStateMachine(t *testing.T) {
	tests := []struct {
		name         string
		pane         ViewPane
		mode         WorkflowMode
		expectedPane ViewPane
		expectedMode WorkflowMode
		expectedLast ViewPane
	}{
		{
			name:         "Success with message pane viewing",
			pane:         MessagePane,
			mode:         Viewing,
			expectedPane: MessagePane,
			expectedMode: Viewing,
			expectedLast: MessagePane,
		},
		{
			name:         "Success with app logs pane fetching",
			pane:         AppLogsPane,
			mode:         Fetching,
			expectedPane: AppLogsPane,
			expectedMode: Fetching,
			expectedLast: AppLogsPane,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(tt.pane, tt.mode)

			assert.Equal(t, tt.expectedPane, sm.ViewPane())
			assert.Equal(t, tt.expectedMode, sm.WorkflowMode())
			assert.Equal(t, tt.expectedLast, sm.LastViewPane())
		})
	}
}

func Test_SetViewPane(t *testing.T) {
	tests := []struct {
		name             string
		initialPane      ViewPane
		setPane          ViewPane
		expectedPane     ViewPane
		expectedLastPane ViewPane
	}{
		{
			name:             "Success when switching to app logs",
			initialPane:      MessagePane,
			setPane:          AppLogsPane,
			expectedPane:     AppLogsPane,
			expectedLastPane: MessagePane,
		},
		{
			name:             "Success when switching to message",
			initialPane:      AppLogsPane,
			setPane:          MessagePane,
			expectedPane:     MessagePane,
			expectedLastPane: MessagePane,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(tt.initialPane, Viewing)
			sm.SetViewPane(tt.setPane)

			assert.Equal(t, tt.expectedPane, sm.ViewPane())
			assert.Equal(t, tt.expectedLastPane, sm.LastViewPane())
		})
	}
}

func Test_EnterViewing(t *testing.T) {
	sm := newStateMachine(MessagePane, Editing)
	sm.EnterViewing(AppLogsPane)

	assert.Equal(t, AppLogsPane, sm.ViewPane())
	assert.Equal(t, Viewing, sm.WorkflowMode())
}

func Test_EnterEditing(t *testing.T) {
	sm := newStateMachine(MessagePane, Viewing)
	sm.EnterEditing()

	assert.Equal(t, Editing, sm.WorkflowMode())
}

func Test_EnterRegenerating(t *testing.T) {
	sm := newStateMachine(MessagePane, Viewing)
	sm.EnterRegenerating()

	assert.Equal(t, Regenerating, sm.WorkflowMode())
}

func Test_EnterFetching(t *testing.T) {
	sm := newStateMachine(MessagePane, Viewing)
	sm.EnterFetching()

	assert.Equal(t, Fetching, sm.WorkflowMode())
}

func Test_IsGenerating(t *testing.T) {
	tests := []struct {
		name     string
		mode     WorkflowMode
		expected bool
	}{
		{
			name:     "Success when fetching generates",
			mode:     Fetching,
			expected: true,
		},
		{
			name:     "Success when regenerating generates",
			mode:     Regenerating,
			expected: true,
		},
		{
			name:     "Success without generating in viewing",
			mode:     Viewing,
			expected: false,
		},
		{
			name:     "Success without generating in editing",
			mode:     Editing,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(MessagePane, tt.mode)
			assert.Equal(t, tt.expected, sm.IsGenerating())
		})
	}
}

func Test_CanEdit(t *testing.T) {
	tests := []struct {
		name     string
		mode     WorkflowMode
		expected bool
	}{
		{
			name:     "Success with editing in viewing",
			mode:     Viewing,
			expected: true,
		},
		{
			name:     "Success without editing in editing",
			mode:     Editing,
			expected: false,
		},
		{
			name:     "Success without editing in fetching",
			mode:     Fetching,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(MessagePane, tt.mode)
			assert.Equal(t, tt.expected, sm.CanEdit())
		})
	}
}

func Test_CanRegenerate(t *testing.T) {
	tests := []struct {
		name     string
		mode     WorkflowMode
		expected bool
	}{
		{
			name:     "Success with regenerate in viewing",
			mode:     Viewing,
			expected: true,
		},
		{
			name:     "Success without regenerate in regenerating",
			mode:     Regenerating,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(MessagePane, tt.mode)
			assert.Equal(t, tt.expected, sm.CanRegenerate())
		})
	}
}

func Test_CanToggleView(t *testing.T) {
	tests := []struct {
		name     string
		mode     WorkflowMode
		expected bool
	}{
		{
			name:     "Success with toggle in viewing",
			mode:     Viewing,
			expected: true,
		},
		{
			name:     "Success without toggle in editing",
			mode:     Editing,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(MessagePane, tt.mode)
			assert.Equal(t, tt.expected, sm.CanToggleView())
		})
	}
}

func Test_CanAccept(t *testing.T) {
	tests := []struct {
		name     string
		mode     WorkflowMode
		expected bool
	}{
		{
			name:     "Success with accept in viewing",
			mode:     Viewing,
			expected: true,
		},
		{
			name:     "Success without accept in editing",
			mode:     Editing,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := newStateMachine(MessagePane, tt.mode)
			assert.Equal(t, tt.expected, sm.CanAccept())
		})
	}
}
