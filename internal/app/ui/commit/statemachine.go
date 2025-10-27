package commit

// ViewPane represents which content pane is currently visible
type ViewPane int

const (
	// MessagePane shows the commit message preview
	MessagePane ViewPane = iota
	// AppLogsPane shows application logs
	AppLogsPane
)

// String returns the string representation of the ViewPane
func (v ViewPane) String() string {
	switch v {
	case MessagePane:
		return "MessagePane"
	case AppLogsPane:
		return "AppLogsPane"
	default:
		return "Unknown"
	}
}

// WorkflowMode represents the current workflow state
type WorkflowMode int

const (
	// Viewing indicates the user is viewing content
	Viewing WorkflowMode = iota
	// Editing indicates the user is editing the commit message
	Editing
	// Fetching indicates the initial commit message is being fetched
	Fetching
	// Regenerating indicates a new commit message is being generated
	Regenerating
)

// String returns the string representation of the WorkflowMode
func (w WorkflowMode) String() string {
	switch w {
	case Viewing:
		return "Viewing"
	case Editing:
		return "Editing"
	case Fetching:
		return "Fetching"
	case Regenerating:
		return "Regenerating"
	default:
		return "Unknown"
	}
}

// stateMachine manages state transitions for the TUI
type stateMachine struct {
	viewPane     ViewPane
	workflowMode WorkflowMode
	lastViewPane ViewPane
}

// newStateMachine creates a new state machine with initial view and mode
func newStateMachine(pane ViewPane, mode WorkflowMode) *stateMachine {
	return &stateMachine{
		viewPane:     pane,
		workflowMode: mode,
		lastViewPane: pane,
	}
}

// ViewPane returns the current view pane
func (sm *stateMachine) ViewPane() ViewPane {
	return sm.viewPane
}

// WorkflowMode returns the current workflow mode
func (sm *stateMachine) WorkflowMode() WorkflowMode {
	return sm.workflowMode
}

// LastViewPane returns the last non-log view pane
func (sm *stateMachine) LastViewPane() ViewPane {
	return sm.lastViewPane
}

// SetViewPane transitions to a new view pane
func (sm *stateMachine) SetViewPane(pane ViewPane) {
	if pane != AppLogsPane {
		sm.lastViewPane = pane
	}
	sm.viewPane = pane
}

// SetWorkflowMode transitions to a new workflow mode
func (sm *stateMachine) SetWorkflowMode(mode WorkflowMode) {
	sm.workflowMode = mode
}

// EnterViewing enters viewing mode at the specified pane
func (sm *stateMachine) EnterViewing(pane ViewPane) {
	sm.SetViewPane(pane)
	sm.SetWorkflowMode(Viewing)
}

// EnterEditing enters editing mode
func (sm *stateMachine) EnterEditing() {
	sm.SetWorkflowMode(Editing)
}

// EnterRegenerating enters regenerating mode
func (sm *stateMachine) EnterRegenerating() {
	sm.SetWorkflowMode(Regenerating)
}

// EnterFetching enters fetching mode
func (sm *stateMachine) EnterFetching() {
	sm.SetWorkflowMode(Fetching)
}

// IsGenerating returns true if the workflow mode is generating
func (sm *stateMachine) IsGenerating() bool {
	return sm.workflowMode == Fetching || sm.workflowMode == Regenerating
}

// CanEdit returns true if the current state allows editing
func (sm *stateMachine) CanEdit() bool {
	return sm.workflowMode == Viewing
}

// CanRegenerate returns true if the current state allows regeneration
func (sm *stateMachine) CanRegenerate() bool {
	return sm.workflowMode == Viewing
}

// CanToggleView returns true if the current state allows view toggling
func (sm *stateMachine) CanToggleView() bool {
	return sm.workflowMode == Viewing
}

// CanAccept returns true if the current state allows accepting the commit
func (sm *stateMachine) CanAccept() bool {
	return sm.workflowMode == Viewing
}
