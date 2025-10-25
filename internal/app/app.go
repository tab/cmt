package app

import (
	"context"
	"os"

	"go.uber.org/fx"

	"cmt/internal/app/cli"
	"cmt/internal/config/logger"
)

// App represents the main application container
type App struct {
	ui  cli.UI
	log logger.Logger
}

// NewApp creates a new application instance with its dependencies
func NewApp(ui cli.UI, log logger.Logger) *App {
	return &App{
		ui:  ui,
		log: log,
	}
}

// Run executes the application with command line arguments
func (a *App) Run() {
	args := os.Args[1:]
	exitCode := a.execute(args)
	os.Exit(exitCode)
}

// execute runs the CLI with given args and handles errors - extracted for testing
func (a *App) execute(args []string) int {
	if err := a.ui.Run(args); err != nil {
		a.log.Error().Err(err).Msg("Application error")
		return 1
	}

	return 0
}

// Register registers the application's lifecycle hooks with fx
func Register(lifecycle fx.Lifecycle, app *App) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
