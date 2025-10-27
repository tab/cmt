package commands

import (
    "context"
    "fmt"

    "cmt/internal/config"
)

// versionCmd is a function-based command for displaying version
type versionCmd func(ctx context.Context, args []string) int

// NewVersionCommand creates a new version command
func NewVersionCommand() Command {
    return versionCmd(func(ctx context.Context, args []string) int {
        fmt.Printf("%s (%s)\n", config.AppName, config.Version)
        return 0
    })
}

// Run executes the version command
func (v versionCmd) Run(ctx context.Context, args []string) int {
    return v(ctx, args)
}
