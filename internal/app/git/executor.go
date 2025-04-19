package git

import (
	"context"
	"os/exec"
)

// Executor represents the git command executor interface
type Executor interface {
	Run(ctx context.Context, name string, arg ...string) *exec.Cmd
}

// executor implements the executor interface
type executor struct{}

// NewGitExecutor creates a new git executor
func NewGitExecutor() Executor {
	return &executor{}
}

// Run executes a git command
func (r *executor) Run(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, arg...)
}
