package gpt

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewGPTClient),
)
