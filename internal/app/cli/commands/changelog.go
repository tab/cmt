package commands

import (
	"context"
	"fmt"
	"strings"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

// changelogCmd handles changelog generation
type changelogCmd struct {
	gitClient git.Client
	gptClient gpt.Client
	log       logger.Logger
}

// NewChangelogCommand creates a new changelog command
func NewChangelogCommand(
	gitClient git.Client,
	gptClient gpt.Client,
	log logger.Logger,
) Command {
	return &changelogCmd{
		gitClient: gitClient,
		gptClient: gptClient,
		log:       log,
	}
}

// Run executes the changelog command
func (c *changelogCmd) Run(ctx context.Context, args []string) int {
	var rangeOpts []string
	for _, arg := range args {
		value := strings.TrimSpace(arg)
		if value != "" {
			rangeOpts = append(rangeOpts, value)
		}
	}

	c.log.Info().
		Str("command", "changelog").
		Str("range", strings.Join(rangeOpts, " ")).
		Msg("Starting changelog generation")

	spin := spinner.New("Fetching git history…")
	spin.Start()

	commits, err := c.gitClient.Log(ctx, rangeOpts)
	if err != nil {
		spin.Stop()
		c.log.Error().
			Str("command", "changelog").
			Err(err).
			Msg("Failed to fetch git log")
		return 1
	}

	spin.SetMessage("Loading…")

	result, err := c.gptClient.FetchChangelog(ctx, commits)
	if err != nil {
		spin.Stop()
		c.log.Error().
			Str("command", "changelog").
			Err(err).
			Msg("Failed to generate changelog")
		return 1
	}

	spin.Stop()

	c.log.Info().
		Str("command", "changelog").
		Int("lines", len(strings.Split(result, "\n"))).
		Msg("Changelog generated successfully")

	fmt.Println(result)
	return 0
}
