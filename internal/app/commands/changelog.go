package commands

import (
	"context"
	"fmt"

	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/app/utils"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

// Changelog defines the interface for changelog command
type Changelog interface {
	Generate(ctx context.Context, args []string) error
}

// changelog represents the changelog command
type changelog struct {
	cfg       *config.Config
	gitClient git.Client
	gptClient gpt.Client
	loader    utils.Loader
	log       logger.Logger
}

// NewChangelog creates a new changelog command instance
func NewChangelog(
	cfg *config.Config,
	gitClient git.Client,
	gptClient gpt.Client,
	loader utils.Loader,
	log logger.Logger,
) Changelog {
	return &changelog{
		cfg:       cfg,
		gitClient: gitClient,
		gptClient: gptClient,
		loader:    loader,
		log:       log,
	}
}

// Generate generates a changelog from git history
func (c *changelog) Generate(ctx context.Context, args []string) error {
	c.loader.Start()

	c.log.Debug().Strs("args", args).Msg("Fetching git log")
	payload, err := c.gitClient.Log(ctx, args)
	if err != nil {
		c.loader.Stop()
		c.log.Debug().Err(err).Msg("Failed to fetch git log")
		errors.HandleGitLogError(err)
		return err
	}

	c.log.Debug().Msg("Requesting changelog from model")
	result, err := c.gptClient.FetchChangelog(ctx, payload)
	if err != nil {
		c.loader.Stop()
		c.log.Debug().Err(err).Msg("Failed to generate changelog")
		errors.HandleModelError(err)
		return err
	}

	c.loader.Stop()

	fmt.Printf("💬 Changelog: \n\n%s\n", result)

	return nil
}
