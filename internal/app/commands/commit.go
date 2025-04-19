package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/app/utils"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

const (
	Accept = "y"
	Edit   = "e"
	Cancel = "n"
)

// Commit defines the interface for commit command
type Commit interface {
	Generate(ctx context.Context, args []string) error
}

// commit represents the commit command
type commit struct {
	cfg       *config.Config
	gitClient git.Client
	gptClient gpt.Client
	loader    utils.Loader
	log       logger.Logger
}

// NewCommit creates a new commit command instance
func NewCommit(
	cfg *config.Config,
	gitClient git.Client,
	gptClient gpt.Client,
	loader utils.Loader,
	log logger.Logger,
) Commit {
	return &commit{
		cfg:       cfg,
		gitClient: gitClient,
		gptClient: gptClient,
		loader:    loader,
		log:       log,
	}
}

// Generate generates a commit message and commits changes
func (c *commit) Generate(ctx context.Context, args []string) error {
	c.loader.Start()

	payload, err := c.gitClient.Diff(ctx)
	if err != nil {
		c.loader.Stop()
		c.log.Debug().Err(err).Msg("Failed to get git diff")
		errors.HandleDiffError(err)
		return err
	}

	c.log.Debug().Msg("Requesting commit message from model")
	result, err := c.gptClient.FetchCommitMessage(ctx, payload)
	if err != nil {
		c.loader.Stop()
		c.log.Debug().Err(err).Msg("Failed to fetch commit message")
		errors.HandleModelError(err)
		return err
	}

	if len(args) > 0 {
		c.log.Debug().Str("prefix", args[0]).Msg("Adding prefix to commit message")
		result = fmt.Sprintf("%s %s", args[0], result)
	}

	c.loader.Stop()

	fmt.Printf("💬 Message: %s", result)
	fmt.Printf("\n\nAccept, edit, or cancel? (%s/%s/%s): ", Accept, Edit, Cancel)

	isAccepted := false

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		c.log.Debug().Err(err).Msg("Failed to read user input")
		errors.HandleInputError(err)
		return err
	}
	answer := strings.TrimSpace(strings.ToLower(input))

	switch answer {
	case Accept:
		c.log.Debug().Msg("User accepted commit message")
		isAccepted = true
	case Edit:
		c.log.Debug().Msg("User requested to edit commit message")
		editedMessage, err := c.gitClient.Edit(context.Background(), result)
		if err != nil {
			c.log.Debug().Err(err).Msg("Failed to edit commit message")
			errors.HandleEditError(err)
			return err
		}

		fmt.Println("\n🧑🏻‍💻 Commit message was changed successfully!")
		result = editedMessage
		isAccepted = true
	default:
		c.log.Debug().Msg("User cancelled commit")
		fmt.Println("❌ Commit aborted")
	}

	if isAccepted {
		ctx, cancel := context.WithTimeout(context.Background(), c.cfg.API.Timeout)
		defer cancel()

		c.log.Debug().Msg("Committing changes")
		output, err := c.gitClient.Commit(ctx, result)
		if err != nil {
			c.log.Debug().Err(err).Msg("Failed to commit changes")
			errors.HandleCommitError(err)
			return err
		}
		fmt.Println("🚀 Changes committed:")
		fmt.Println(output)
	}

	return nil
}
