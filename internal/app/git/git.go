package git

import (
	"bytes"
	"context"
	"os"
	"strings"

	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

// Client represents a git client interface
type Client interface {
	Diff(ctx context.Context) (string, error)
	Status(ctx context.Context) (string, error)
	Log(ctx context.Context, opts []string) (string, error)
	Commit(ctx context.Context, message string) (string, error)
}

// client implements the git client interface
type client struct {
	cfg      *config.Config
	executor Executor
	log      logger.Logger
}

// NewGitClient creates a new git client
func NewGitClient(
	cfg *config.Config,
	executor Executor,
	log logger.Logger,
) Client {
	return &client{
		cfg:      cfg,
		executor: executor,
		log:      log,
	}
}

// Diff returns the git diff
func (g *client) Diff(ctx context.Context) (string, error) {
	args := []string{"diff", "--staged", "--minimal", "--ignore-all-space", "--ignore-blank-lines"}

	g.log.Debug().Strs("args", args).Msg("Running git diff command")
	cmd := g.executor.Run(ctx, "git", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		g.log.Error().Err(err).Msg("Failed to execute git diff command")
		return "", errors.ErrFailedToLoadGitDiff
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", errors.ErrNoGitChanges
	}

	g.log.Debug().Int("diff_length", len(result)).Msg("Git diff loaded successfully")
	return result, nil
}

// Status returns the git status for staged files
func (g *client) Status(ctx context.Context) (string, error) {
	args := []string{"diff", "--staged", "--name-status"}

	g.log.Debug().Strs("args", args).Msg("Running git status command")
	cmd := g.executor.Run(ctx, "git", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		g.log.Error().Err(err).Msg("Failed to execute git status command")
		return "", errors.ErrFailedToLoadGitDiff
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", errors.ErrNoGitChanges
	}

	g.log.Debug().Int("status_length", len(result)).Msg("Git status loaded successfully")
	return result, nil
}

// Log returns the git log with detailed format: hash|subject|author|date
func (g *client) Log(ctx context.Context, opts []string) (string, error) {
	args := []string{"log", "--format=%h|%s|%an|%ar"}
	args = append(args, opts...)

	g.log.Debug().Strs("args", args).Msg("Running git log command")
	cmd := g.executor.Run(ctx, "git", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		g.log.Error().Err(err).Msg("Failed to execute git log command")
		return "", errors.ErrFailedToLoadGitLog
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		g.log.Info().Msg("No git commits found")
		return "", errors.ErrNoGitCommits
	}

	g.log.Debug().Int("log_length", len(result)).Msg("Git log loaded successfully")
	return result, nil
}

// Commit commits the staged git changes
func (g *client) Commit(ctx context.Context, message string) (string, error) {
	if message == "" {
		g.log.Error().Msg("Commit message is empty")
		return "", errors.ErrCommitMessageEmpty
	}

	g.log.Debug().Msg("Running git commit command")
	cmd := g.executor.Run(ctx, "git", "commit", "-m", message)

	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		g.log.Error().Err(err).Msg("Failed to execute git commit command")
		return "", errors.ErrFailedToCommit
	}

	result := out.String()
	g.log.Debug().Msg("Successfully committed changes")
	return result, nil
}
