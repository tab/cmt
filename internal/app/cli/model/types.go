package model

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"

	"cmt/internal/app/cli/components"
	"cmt/internal/app/cli/workflow"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

//go:generate mockgen -source=types.go -destination=types_mock.go -package=model

type workflowService interface {
	GenerateCommit(ctx context.Context, prefix string) (workflow.CommitResult, error)
	GenerateChangelog(ctx context.Context, between string) (workflow.ChangelogResult, error)
}

// Action represents what the user wants to do next
type Action int

const (
	ActionNone Action = iota
	ActionAccept
	ActionEdit
	ActionRefresh
	ActionCancel
)

// Mode represents the vim editing mode
type Mode int

const (
	EditModeNormal Mode = iota
	EditModeInsert
)

// Flow represents the user's workflow intention
type Flow int

const (
	FlowCommit Flow = iota
)

// State represents the current state of the application
type State int

const (
	StateInit State = iota
	StateFetch
	StateViewCommit
	StateEditCommit
	StateViewLogs
	StateExit
)

// Model is the Bubble Tea model
type Model struct {
	Ctx context.Context

	Cfg      *config.Config
	Workflow workflowService

	controllers map[State]stateController

	State State

	Spinner       spinner.Model
	Textarea      textarea.Model
	Viewport      viewport.Model
	LogBuffer     *logger.LogBuffer
	PreviousState State

	Content  string
	FileTree *components.FileTree
	Err      error
	Width    int
	Height   int

	UserFlow   Flow
	UserAction Action

	Mode Mode

	// FetchStarted prevents duplicate fetch operations when spinner ticks
	// before the initial async fetch completes.
	FetchStarted bool

	Prefix string

	Log logger.Logger
}

// Result is a unified message for all async operations
type Result struct {
	Content  string
	FileTree *components.FileTree
	Err      error
}
