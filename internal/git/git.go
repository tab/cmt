package git

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"cmt/internal/errors"
)

// Executor the executor interface
type Executor interface {
	Run(ctx context.Context, name string, arg ...string) *exec.Cmd
}

// GitClient the git client interface
type GitClient interface {
	Diff(ctx context.Context, opts []string) (string, error)
	Log(ctx context.Context, opts []string) (string, error)
	Edit(ctx context.Context, message string) (string, error)
	Commit(ctx context.Context, message string) (string, error)
}

// NewGitClient creates a new git client instance
func NewGitClient() GitClient {
	return &Git{
		Executor: &GitExecutor{},
	}
}

// Git represents the git client
type Git struct {
	Executor Executor
}

// GitExecutor the git executor
type GitExecutor struct{}

// Run runs a git command
func (r *GitExecutor) Run(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, arg...)
}

// Diff returns the git diff
func (g *Git) Diff(ctx context.Context, opts []string) (string, error) {
	args := []string{"diff", "--minimal", "--ignore-all-space", "--ignore-blank-lines"}

	if len(opts) == 0 {
		args = append(args, "--staged")
	} else {
		args = append(args, opts...)
	}

	cmd := g.Executor.Run(ctx, "git", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.ErrFailedToLoadGitDiff
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", errors.ErrNoGitChanges
	}

	return result, nil
}

// Log returns the git log
func (g *Git) Log(ctx context.Context, opts []string) (string, error) {
	args := []string{"log", "--format='%h %s %b'"}
	args = append(args, opts...)

	cmd := g.Executor.Run(ctx, "git", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.ErrFailedToLoadGitLog
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", errors.ErrNoGitCommits
	}

	return result, nil
}

// Edit edits the commit message
func (g *Git) Edit(ctx context.Context, message string) (string, error) {
	tmpFile, err := os.CreateTemp("", "editor")
	if err != nil {
		return "", errors.ErrFailedToCreateFile
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(message)
	if err != nil {
		return "", errors.ErrFailedToWriteFile
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := g.Executor.Run(ctx, editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", errors.ErrFailedToRunEditor
	}

	editedMessageBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", errors.ErrFailedToReadFile
	}
	message = strings.TrimSpace(string(editedMessageBytes))

	return message, nil
}

// Commit commits the staged git changes
func (g *Git) Commit(ctx context.Context, message string) (string, error) {
	if message == "" {
		return "", errors.ErrCommitMessageEmpty
	}

	cmd := g.Executor.Run(ctx, "git", "commit", "-m", message)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.ErrFailedToCommit
	}

	result := out.String()
	return result, nil
}
