package cli

import (
	"go.uber.org/fx"

	"cmt/internal/app/cli/commands"
	"cmt/internal/app/cli/spinner"
)

// Module exports the CLI module for dependency injection
var Module = fx.Options(
	commands.Module,
	spinner.Module,
	fx.Provide(
		NewRunner,
		NewCLI,
	),
)
