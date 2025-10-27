package spinner

import "go.uber.org/fx"

// Factory is a function type that creates new spinner instances
type Factory func() Model

// Module exports the spinner module for dependency injection
var Module = fx.Options(
	fx.Provide(func() Factory {
		return NewSpinner
	}),
)
