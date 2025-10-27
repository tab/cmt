package app

import (
	"go.uber.org/fx"

	"cmt/internal/app/cli"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

var Module = fx.Options(
	cli.Module,
	git.Module,
	gpt.Module,
	logger.Module,
)
