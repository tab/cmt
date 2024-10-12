package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"cmt/internal/errors"
)

type Executor interface {
	Run(ctx context.Context, name string, arg ...string) *exec.Cmd
}

type Git struct {
	Executor Executor
}

type GitExecutor struct{}

func (r *GitExecutor) Run(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, arg...)
}

func (g *Git) Diff(ctx context.Context) (string, error) {
	cmd := g.Executor.Run(ctx, "git", "diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff error: %w", err)
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", errors.ErrNoGitChanges
	}

	return result, nil
}

func (g *Git) Commit(ctx context.Context, message string) (string, error) {
	if message == "" {
		return "", errors.ErrCommitMessageEmpty
	}

	cmd := g.Executor.Run(ctx, "git", "commit", "-m", message)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git commit error: %w", err)
	}

	result := out.String()
	return result, nil
}
