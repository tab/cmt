package git

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewGitExecutor,
		NewGitClient,
	),
)
