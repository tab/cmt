package app

import (
	"context"
	"os"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"cmt/internal/app/cli"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

// Run is the main entry point for the application
func Run() int {
	cfg, err := config.Load()
	if err != nil {
		return 1
	}

	ctx := context.Background()

	fxApp, exitCode := createFxApp(ctx, cfg, Module)

	if err := fxApp.Start(ctx); err != nil {
		return 1
	}

	if err := fxApp.Stop(ctx); err != nil {
		return 1
	}

	return exitCode
}

// createFxApp creates and configures the FX application
func createFxApp(ctx context.Context, cfg *config.Config, module fx.Option) (*fx.App, int) {
	args := os.Args[1:]

	var exitCode int

	return fx.New(
		fx.WithLogger(createFxLogger(cfg)),
		fx.Supply(cfg, args),
		module,
		fx.Invoke(func(cliInstance *cli.CLI) {
			exitCode = cliInstance.Run(ctx, args)
		}),
	), exitCode
}

// createFxLogger returns an FX logger based on the config
func createFxLogger(cfg *config.Config) func() fxevent.Logger {
	return func() fxevent.Logger {
		if strings.EqualFold(cfg.Logging.Level, logger.DebugLevel) {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		}
		return fxevent.NopLogger
	}
}
