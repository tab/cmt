package commit

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cmt/internal/commands"
	"cmt/internal/errors"
	"cmt/internal/utils"
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
	fmt.Print("\n\nAccept? (y/n): ")

	var answer string
	if c.InputReader != nil {
		answer, err = c.InputReader()
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
		answer = strings.TrimSpace(input)
	}

	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" {
		output, err := c.Options.Client.Commit(c.Options.Ctx, commitMessage)
		if err != nil {
			errors.HandleCommitError(err)
			return err
		}
		fmt.Println("🚀 Changes committed:")
		fmt.Println(output)
	} else {
		fmt.Println("❌ Commit aborted")
	}

	return nil
}
