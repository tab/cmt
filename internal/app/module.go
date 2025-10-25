package app

import (
	"context"

	"go.uber.org/fx"

	"cmt/internal/app/cli"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

// NewContext provides a background context for the application
func NewContext() context.Context {
	return context.Background()
}

var Module = fx.Options(
	fx.Provide(NewContext),
	cli.Module,
	git.Module,
	gpt.Module,
	logger.Module,
	fx.Provide(NewApp),
	fx.Invoke(Register),
)
