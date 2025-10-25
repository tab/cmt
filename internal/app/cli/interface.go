package cli

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// UI defines the interface for terminal user interface operations
type UI interface {
	Run(args []string) error
}

// Program defines the interface for Bubble Tea program operations
type Program interface {
	Run() (tea.Model, error)
}

// Builder creates a Bubble Tea program from a model
type Builder func(model tea.Model, opts ...tea.ProgramOption) Program

// gitClient defines git operations
type gitClient interface {
	Commit(ctx context.Context, msg string) (string, error)
}
