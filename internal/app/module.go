package app

import (
	"go.uber.org/fx"

	"cmt/internal/app/cli"
	"cmt/internal/app/commands"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/app/utils"
	"cmt/internal/config/logger"
)

var Module = fx.Options(
	cli.Module,
	commands.Module,
	git.Module,
	gpt.Module,
	logger.Module,
	utils.Module,
	fx.Provide(NewApp),
	fx.Invoke(Register),
)
