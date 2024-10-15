package commands

import (
	"context"

	"cmt/internal/errors"
	"cmt/internal/git"
	"cmt/internal/gpt"
)

type GenerateOptions struct {
	Ctx    context.Context
	Client git.GitClient
	Model  gpt.GPTModelClient
	Args   []string
}

type Command interface {
	Generate() error
}

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
