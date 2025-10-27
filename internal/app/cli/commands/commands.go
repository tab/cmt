package commands

import (
	"context"

	"go.uber.org/fx"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

// Command represents a CLI command interface
type Command interface {
	Run(ctx context.Context, args []string) int
}

// CommandsParams groups dependencies for constructing commands
type CommandsParams struct {
	fx.In

	GitClient git.Client
	GPTClient gpt.Client
	Log       logger.Logger
	Spinner   spinner.Factory
}

// commandsResult groups all commands for FX injection
type commandsResult struct {
	fx.Out

	Help      Command `name:"help"`
	Version   Command `name:"version"`
	Changelog Command `name:"changelog"`
	Commit    Command `name:"commit"`
}

// provideCommands creates all command instances
func provideCommands(p CommandsParams) commandsResult {
	return commandsResult{
		Help:      NewHelpCommand(),
		Version:   NewVersionCommand(),
		Changelog: NewChangelogCommand(p.GitClient, p.GPTClient, p.Log),
		Commit:    NewCommitCommand(p.GitClient, p.GPTClient, p.Log, p.Spinner),
	}
}
