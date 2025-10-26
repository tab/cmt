package cli

import (
	"strings"

	"go.uber.org/fx"

	"cmt/internal/app/cli/commands"
	"cmt/internal/app/errors"
)

// Runner represents a command resolver interface
type Runner interface {
	Resolve(args []string) (commands.Command, []string, error)
}

// Params holds dependencies for NewRunner constructor
type Params struct {
	fx.In

	Help      commands.Command `name:"help"`
	Version   commands.Command `name:"version"`
	Changelog commands.Command `name:"changelog"`
	Commit    commands.Command `name:"commit"`
}

// runner implements the Runner interface
type runner struct {
	help        commands.Command
	version     commands.Command
	changelog   commands.Command
	commit      commands.Command
	dispatchMap map[string]commands.Command
}

// NewRunner creates a new command runner with a pre-built dispatch map
func NewRunner(p Params) Runner {
	r := &runner{
		help:        p.Help,
		version:     p.Version,
		changelog:   p.Changelog,
		commit:      p.Commit,
		dispatchMap: make(map[string]commands.Command),
	}

	r.dispatchMap["changelog"] = r.changelog
	r.dispatchMap["--changelog"] = r.changelog
	r.dispatchMap["-c"] = r.changelog

	r.dispatchMap["help"] = r.help
	r.dispatchMap["--help"] = r.help
	r.dispatchMap["-h"] = r.help

	r.dispatchMap["version"] = r.version
	r.dispatchMap["--version"] = r.version
	r.dispatchMap["-v"] = r.version

	return r
}

// Resolve parses arguments and returns the appropriate command
func (r *runner) Resolve(args []string) (commands.Command, []string, error) {
	if len(args) == 0 {
		return r.commit, []string{}, nil
	}

	command := strings.ToLower(args[0])

	if cmd, ok := r.dispatchMap[command]; ok {
		remainingArgs := []string{}
		if command == "changelog" || command == "--changelog" || command == "-c" {
			remainingArgs = args[1:]
		}
		return cmd, remainingArgs, nil
	}

	switch {
	case strings.HasPrefix(command, "--prefix") || command == "prefix" || command == "-p":
		return r.commit, args, nil
	default:
		return nil, nil, errors.ErrUnknownCommand
	}
}
