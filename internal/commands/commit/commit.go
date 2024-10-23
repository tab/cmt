package commit

import (
	"fmt"
	"strings"

	"cmt/internal/commands"
	"cmt/internal/errors"
	"cmt/internal/utils"
)

const (
	Accept = "y"
	Edit   = "e"
	Cancel = "n"
)

type Command struct {
	Options     commands.GenerateOptions
	InputReader func() (string, error)
}

func NewCommand(options commands.GenerateOptions, inputReader func() (string, error)) *Command {
	return &Command{
		Options:     options,
		InputReader: inputReader,
	}
}

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
		editedMessage, err := c.Options.Client.Edit(c.Options.Ctx, commitMessage)
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
