package workflow

import (
	"context"
	"fmt"

	"cmt/internal/app/cli/components"
	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

// Service exposes high-level workflows for generating commit messages and changelogs.
type Service interface {
	GenerateCommit(ctx context.Context, prefix string) (CommitResult, error)
	GenerateChangelog(ctx context.Context, between string) (ChangelogResult, error)
}

type service struct {
	git git.Client
	gpt gpt.Client
	log logger.Logger
}

// CommitResult carries the artifacts required to review a generated commit.
type CommitResult struct {
	Message  string
	FileTree *components.FileTree
}

// ChangelogResult carries the data required to review a generated changelog.
type ChangelogResult struct {
	Content string
}

// NewService wires the workflow service with git, GPT and logging dependencies.
func NewService(gitClient git.Client, gptClient gpt.Client, log logger.Logger) Service {
	return &service{
		git: gitClient,
		gpt: gptClient,
		log: log,
	}
}

// GenerateCommit orchestrates the end-to-end flow for producing a commit message.
func (s *service) GenerateCommit(ctx context.Context, prefix string) (CommitResult, error) {
	s.log.Info().Msg("Fetching git diff")
	diff, err := s.git.Diff(ctx)
	if err != nil {
		return CommitResult{}, err
	}

	s.log.Info().Msg("Fetching git status for file list")
	status, err := s.git.Status(ctx)
	var fileTree *components.FileTree
	if err != nil {
		s.log.Warn().Err(err).Msg("Failed to fetch git status, continuing without file list")
	} else {
		fileTree = components.ParseGitStatus(status)
	}

	s.log.Info().Msg("Generating commit message")
	message, err := s.gpt.FetchCommitMessage(ctx, diff)
	if err != nil {
		return CommitResult{}, err
	}

	if prefix != "" {
		s.log.Info().Str("prefix", prefix).Msg("Applied prefix to commit message")
		message = fmt.Sprintf("%s %s", prefix, message)
	}

	s.log.Info().Msg("Commit message generated successfully")

	return CommitResult{
		Message:  message,
		FileTree: fileTree,
	}, nil
}

// GenerateChangelog orchestrates the changelog generation workflow.
func (s *service) GenerateChangelog(ctx context.Context, between string) (ChangelogResult, error) {
	opts := []string{}
	logEvent := s.log.Info()
	if between != "" {
		opts = append(opts, between)
		logEvent.Str("range", between).Msg("Fetching git log for changelog")
	} else {
		logEvent.Msg("Fetching git log for changelog")
	}

	logOutput, err := s.git.Log(ctx, opts)
	if err != nil {
		return ChangelogResult{}, err
	}

	if logOutput == "" {
		return ChangelogResult{}, errors.ErrNoGitCommits
	}

	s.log.Info().Msg("Generating changelog")
	payload := components.ConvertCommitLogForGPT(logOutput)
	changelog, err := s.gpt.FetchChangelog(ctx, payload)
	if err != nil {
		return ChangelogResult{}, err
	}

	s.log.Info().Msg("Changelog generated successfully")

	return ChangelogResult{
		Content: changelog,
	}, nil
}
