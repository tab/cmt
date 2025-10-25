package workflow

import "go.uber.org/fx"

// Module exposes the workflow service to the fx graph.
var Module = fx.Options(
	fx.Provide(NewService),
)
