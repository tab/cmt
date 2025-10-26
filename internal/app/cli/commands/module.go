package commands

import (
    "go.uber.org/fx"
)

// Module exports the commands module for dependency injection
var Module = fx.Options(
    fx.Provide(provideCommands),
)
