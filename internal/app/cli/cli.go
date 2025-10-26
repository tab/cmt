package cli

import (
	"context"
	"fmt"

	"cmt/internal/app/errors"
)

// CLI represents the command line interface handler
type CLI struct {
	runner Runner
}

// NewCLI creates a new CLI instance
func NewCLI(runner Runner) *CLI {
	return &CLI{
		runner: runner,
	}
}

// Run processes command line arguments and executes the appropriate command
func (c *CLI) Run(ctx context.Context, opts []string) int {
	cmd, args, err := c.runner.Resolve(opts)

	if err != nil {
		fmt.Println(errors.Format(err))
		return 1
	}

	return cmd.Run(ctx, args)
}
