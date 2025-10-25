package main

import (
	"context"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"cmt/internal/app"
	"cmt/internal/app/changelog"
	"cmt/internal/app/cli/helpers"
	"cmt/internal/app/cli/workflow"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

// main is the entry point for the application
func main() {
	runApp()
}

// runApp contains the main application logic
func runApp() {
	args := os.Args[1:]

	// Handle simple commands without FX
	if handleSimpleCommands(args) {
		return
	}

	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}

	// Handle changelog command with FX
	if handleChangelogCommand(args, cfg) {
		return
	}

	// Handle commit flow with FX
	application := createApp(cfg)
	application.Run()
}

// createApp creates the FX application with the given config
func createApp(cfg *config.Config) *fx.App {
	return fx.New(
		fx.WithLogger(createFxLogger(cfg)),
		fx.Supply(cfg),
		app.Module,
	)
}

// createFxLogger returns an FX logger based on the config
func createFxLogger(cfg *config.Config) func() fxevent.Logger {
	return func() fxevent.Logger {
		if cfg.Logging.Level == logger.DebugLevel {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		}
		return fxevent.NopLogger
	}
}

// handleSimpleCommands handles commands that don't need FX (help, version)
func handleSimpleCommands(args []string) bool {
	if len(args) == 0 {
		return false
	}

	cmd := args[0]
	switch {
	case helpers.IsHelpCmd(cmd):
		helpers.RenderHelp()
		return true
	case helpers.IsVersionCmd(cmd):
		helpers.RenderVersion()
		return true
	default:
		return false
	}
}

// handleChangelogCommand handles changelog command with FX DI
func handleChangelogCommand(args []string, cfg *config.Config) bool {
	if len(args) == 0 {
		return false
	}

	cmd := args[0]
	if !helpers.IsChangelogCmd(cmd) {
		return false
	}

	between := ""
	if len(args) > 1 {
		between = args[1]
	}

	changelogApp := createChangelogApp(cfg, between)
	changelogApp.Run()
	return true
}

// createChangelogApp creates an FX app for changelog generation
func createChangelogApp(cfg *config.Config, between string) *fx.App {
	return fx.New(
		fx.WithLogger(createFxLogger(cfg)),
		fx.Supply(cfg),
		fx.Provide(app.NewContext),
		logger.Module,
		git.Module,
		gpt.Module,
		workflow.Module,
		changelog.Module,
		fx.Invoke(func(lc fx.Lifecycle, gen *changelog.Generator, shutdowner fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := gen.Generate(between); err != nil {
							os.Exit(1)
						}
						shutdowner.Shutdown()
					}()
					return nil
				},
			})
		}),
	)
}
