package commands

import (
	"context"
	"fmt"

	"cmt/internal/config"
)

// helpCmd is a function-based command for displaying help
type helpCmd func(ctx context.Context, args []string) int

// NewHelpCommand creates a new help command
func NewHelpCommand() Command {
	return helpCmd(func(ctx context.Context, args []string) int {
		fmt.Print(GetUsage())
		return 0
	})
}

// Run executes the help command
func (h helpCmd) Run(ctx context.Context, args []string) int {
	return h(ctx, args)
}

// GetUsage returns the usage text for the CLI
func GetUsage() string {
	return fmt.Sprintf(`%s (%s) - %s

Usage:
  cmt [command] [options]

Commands:
  changelog [RANGE]   Generate a changelog from git history
  version             Display version information
  help                Display this help message

Examples:
  cmt                        Generate commit message for staged changes
  cmt --prefix "TASK-123"     Add "TASK-123" prefix to commit message
  cmt changelog              Generate changelog for all commits
  cmt changelog v1.0..v2.0   Generate changelog between versions
  cmt --version              Show version
  cmt --help                 Show this help

Navigation:
  tab                 Switch focus between panes
  j/k, ↑/↓            Scroll focused pane
  a                   Accept and commit
  e                   Edit commit message
  r                   Regenerate commit message
  l                   Toggle application logs
  q, Ctrl+C           Quit without committing

Environment:
  OPENAI_API_KEY         Required: Your OpenAI API key
`,
		config.AppName,
		config.Version,
		config.AppDescription,
	)
}
