package changelog

import (
	"go.uber.org/fx"
)

// Module provides changelog generation dependencies
var Module = fx.Module(
	"changelog",
	fx.Provide(NewGenerator),
)
