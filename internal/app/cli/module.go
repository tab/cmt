package cli

import (
	"go.uber.org/fx"

	"cmt/internal/app/cli/workflow"
)

var Module = fx.Options(
	workflow.Module,
	fx.Provide(NewUI),
)
