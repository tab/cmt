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
	Log(ctx context.Context, opts []string) (string, error)
	Edit(ctx context.Context, message string) (string, error)
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

// Log returns the git log
func (g *client) Log(ctx context.Context, opts []string) (string, error) {
	args := []string{"log", "--format='%h %s %b'"}
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

// Edit edits the commit message
func (g *client) Edit(ctx context.Context, message string) (string, error) {
	tmpFile, err := os.CreateTemp("", "editor")
	if err != nil {
		g.log.Error().Err(err).Msg("Failed to create temp file")
		return "", errors.ErrFailedToCreateFile
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(message)
	if err != nil {
		g.log.Error().Err(err).Msg("Failed to write to temp file")
		return "", errors.ErrFailedToWriteFile
	}
	tmpFile.Close()

	editor := g.cfg.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
	}

	g.log.Debug().Str("editor", editor).Str("file", tmpFile.Name()).Msg("Opening editor")
	cmd := g.executor.Run(ctx, editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		g.log.Error().Err(err).Str("editor", editor).Msg("Failed to run editor")
		return "", errors.ErrFailedToRunEditor
	}

	editedMessageBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		g.log.Error().Err(err).Msg("Failed to read edited file")
		return "", errors.ErrFailedToReadFile
	}
	message = strings.TrimSpace(string(editedMessageBytes))

	g.log.Debug().Int("message_length", len(message)).Msg("Successfully edited message")
	return message, nil
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
