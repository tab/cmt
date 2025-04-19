package cli

import (
	"context"
	"fmt"
	"os"

	"cmt/internal/app/commands"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

const (
	Usage = `Usage:
  cmt                            Generate a commit message based on staged changes
  cmt prefix <type>              Generate a commit message with a prefix
  cmt changelog [commit range]   Generate a changelog for a range of commits
  cmt help                       Show help
  cmt version                    Show version

Examples:
  cmt prefix TASK-1234           Generate a commit message with a task prefix
  cmt changelog                  From initial commit to HEAD
  cmt changelog v1.0.0..v1.1.0   From v1.0.0 to v1.1.0
  cmt changelog 2606b09..5e3ac73 From 2606b09 to 5e3ac73
  cmt changelog 2606b09..HEAD    From 2606b09 to HEAD`
)

// CLI defines the interface for cli operations
type CLI interface {
	Run(args []string) error
}

// cli represents the command-line interface for the application
type cli struct {
	commit    commands.Commit
	changelog commands.Changelog
	log       logger.Logger
}

// NewCLI creates a new cli instance
func NewCLI(
	commit commands.Commit,
	changelog commands.Changelog,
	log logger.Logger,
) CLI {
	return &cli{
		commit:    commit,
		changelog: changelog,
		log:       log,
	}
}

// Run processes command-line arguments and executes commands
func (c *cli) Run(args []string) error {
	if len(args) == 0 {
		c.handleCommit(args)

		os.Exit(0)
		return nil
	}

	cmd := args[0]
	params := args[1:]

	switch cmd {
	case "prefix", "--prefix", "-p":
		c.handleCommit(params)
	case "changelog":
		c.handleChangelog(params)
	case "help", "--help", "-h":
		c.handleHelp()
	case "version", "--version", "-v":
		c.handleVersion()
	default:
		c.handleUnknown()
	}

	os.Exit(0)
	return nil
}

// handleHelp displays help information
func (c *cli) handleHelp() {
	c.log.Debug().Msg("Displaying help information")
	fmt.Println(Usage)
}

// handleVersion displays version information
func (c *cli) handleVersion() {
	c.log.Debug().Msg("Displaying version information")
	fmt.Printf("Version: %s\n", config.Version)
}

// handleCommit generates a commit message based on the provided parameters
func (c *cli) handleCommit(params []string) {
	c.log.Debug().Msg("Generating commit message")
	_ = c.commit.Generate(context.Background(), params)
}

// handleChangelog generates a changelog based on the provided parameters
func (c *cli) handleChangelog(params []string) {
	c.log.Debug().Msg("Generating changelog")
	_ = c.changelog.Generate(context.Background(), params)
}

// handleUnknown handles unknown commands
func (c *cli) handleUnknown() {
	c.log.Debug().Msg("Unknown command")
	fmt.Println("Unknown command. Use 'cmt help' for more information")
}
