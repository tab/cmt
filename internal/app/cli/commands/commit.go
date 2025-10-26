package commands

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/app/ui/commit"
	"cmt/internal/config/logger"
)

// commitCmd handles commit message generation and committing
type commitCmd struct {
	gitClient git.Client
	gptClient gpt.Client
	log       logger.Logger
	spinner   spinner.Factory
}

// NewCommitCommand creates a new commit command
func NewCommitCommand(
	gitClient git.Client,
	gptClient gpt.Client,
	log logger.Logger,
	spinner spinner.Factory,
) Command {
	return &commitCmd{
		gitClient: gitClient,
		gptClient: gptClient,
		spinner:   spinner,
		log:       log,
	}
}

// Run executes the commit command
func (c *commitCmd) Run(ctx context.Context, args []string) int {
	prefix := c.parsePrefix(args)

	c.log.Info().
		Str("command", "commit").
		Str("prefix", prefix).
		Msg("Starting commit workflow")

	c.log.Info().
		Str("command", "commit").
		Msg("Launching TUI")

	input := commit.Input{
		Ctx:       ctx,
		Prefix:    prefix,
		GitClient: c.gitClient,
		GPTClient: c.gptClient,
		Spinner:   c.spinner,
		Logger:    c.log,
	}

	model := commit.NewModel(input)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("TUI error: %v\n", err)
		return 1
	}

	commitModel, ok := finalModel.(commit.Model)
	if !ok {
		return 1
	}

	output := commitModel.GetOutput()

	if output.Error != nil {
		fmt.Printf("Commit workflow failed: %v\n", output.Error)
		return 1
	}

	if !output.Accepted {
		fmt.Println("❌ Commit cancelled")
		return 0
	}

	fmt.Println("🚀 Changes committed:")
	if output.CommitOutput != "" {
		fmt.Println(output.CommitOutput)
	}
	return 0
}

// parsePrefix extracts prefix from command arguments
func (c *commitCmd) parsePrefix(args []string) string {
	for i, arg := range args {
		lowerArg := strings.ToLower(arg)

		if strings.HasPrefix(lowerArg, "--prefix=") {
			value := strings.TrimPrefix(arg, "--prefix=")
			value = strings.TrimPrefix(value, "--PREFIX=")
			return strings.TrimSpace(value)
		}

		if strings.HasPrefix(lowerArg, "-p=") {
			value := strings.TrimPrefix(arg, "-p=")
			value = strings.TrimPrefix(value, "-P=")
			return strings.TrimSpace(value)
		}

		if lowerArg == "--prefix" || lowerArg == "-p" || lowerArg == "prefix" {
			if i+1 < len(args) {
				nextArg := args[i+1]
				if !strings.HasPrefix(nextArg, "-") {
					return strings.TrimSpace(nextArg)
				}
			}
		}
	}
	return ""
}
