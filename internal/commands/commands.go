package commands

import (
	"context"

	"cmt/internal/errors"
	"cmt/internal/git"
	"cmt/internal/gpt"
)

// GenerateOptions represents the options for the Generate command
type GenerateOptions struct {
	Ctx    context.Context
	Client git.GitClient
	Model  gpt.GPTModelClient
	Args   []string
}

// Command represents the command interface
type Command interface {
	Generate() error
}

// ValidateOptions validates the generate command options
func ValidateOptions(opts GenerateOptions) error {
	if opts.Ctx == nil {
		return errors.ErrInvalidContext
	}

	if opts.Client == nil {
		return errors.ErrNilClient
	}

	if opts.Model == nil {
		return errors.ErrNilModel
	}

	return nil
}
