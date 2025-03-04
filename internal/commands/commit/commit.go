package commit

import (
	"context"
	"fmt"
	"strings"

	"cmt/internal/commands"
	"cmt/internal/config"
	"cmt/internal/errors"
	"cmt/internal/utils"
)

const (
	Accept = "y"
	Edit   = "e"
	Cancel = "n"
)

// Command represents the commit command
type Command struct {
	Options     commands.GenerateOptions
	InputReader func() (string, error)
}

// NewCommand creates a new command instance
func NewCommand(options commands.GenerateOptions, inputReader func() (string, error)) *Command {
	return &Command{
		Options:     options,
		InputReader: inputReader,
	}
}

// Generate generates a commit message
func (c *Command) Generate() error {
	loader := utils.NewLoader()
	loader.Start()

	diff, err := c.Options.Client.Diff(c.Options.Ctx, nil)
	if err != nil {
		loader.Stop()
		errors.HandleDiffError(err)
		return err
	}

	commitMessage, err := c.Options.Model.FetchCommitMessage(c.Options.Ctx, diff)
	loader.Stop()

	if err != nil {
		return fmt.Errorf("error requesting commit message: %w", err)
	}

	if len(c.Options.Args) > 0 {
		commitMessage = fmt.Sprintf("%s %s", c.Options.Args[0], commitMessage)
	}

	fmt.Printf("💬 Message: %s", commitMessage)
	fmt.Printf("\n\nAccept, edit, or cancel? (%s/%s/%s): ", Accept, Edit, Cancel)

	isAccepted := false

	answer, err := c.InputReader()
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))

	switch answer {
	case Accept:
		isAccepted = true
	case Edit:
		editedMessage, err := c.Options.Client.Edit(context.Background(), commitMessage)
		if err != nil {
			errors.HandleEditError(err)
			return err
		}

		fmt.Println("\n🧑🏻‍💻 Commit message was changed successfully!")
		commitMessage = editedMessage
		isAccepted = true
	default:
		fmt.Println("❌ Commit aborted")
	}

	if isAccepted {
		ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
		c.Options.Ctx = ctx
		defer cancel()

		output, err := c.Options.Client.Commit(c.Options.Ctx, commitMessage)
		if err != nil {
			errors.HandleCommitError(err)
			return err
		}
		fmt.Println("🚀 Changes committed:")
		fmt.Println(output)
	}

	return nil
}
