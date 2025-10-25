package cli

import (
	"context"
	"fmt"

	"cmt/internal/app/cli/model"
	"cmt/internal/app/cli/workflow"

	tea "github.com/charmbracelet/bubbletea"

	"cmt/internal/app/errors"
	"cmt/internal/app/git"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

type ui struct {
	ctx      context.Context
	cfg      *config.Config
	builder  Builder
	git      gitClient
	workflow workflow.Service
	log      logger.Logger
}

// NewUI creates a new TUI instance
func NewUI(
	ctx context.Context,
	cfg *config.Config,
	gitClient git.Client,
	workflowService workflow.Service,
	log logger.Logger,
) UI {
	return &ui{
		ctx:      ctx,
		cfg:      cfg,
		builder:  newProgramBuilder,
		git:      gitClient,
		workflow: workflowService,
		log:      log,
	}
}

// newProgramBuilder is an adapter that wraps tea.NewProgram to match Builder signature
func newProgramBuilder(model tea.Model, opts ...tea.ProgramOption) Program {
	return tea.NewProgram(model, opts...)
}

// Run starts the TUI application
func (t *ui) Run(args []string) error {
	t.log.Info().Msg("Starting application UI")
	m := model.New(t.ctx, t.cfg, t.workflow, t.log, args)

	p := t.builder(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		t.log.Error().Err(err).Msg("Bubble Tea program error")
		return fmt.Errorf("TUI error: %w", err)
	}

	final, ok := finalModel.(model.Model)
	if !ok {
		t.log.Error().Msg("Unexpected model type returned from Bubble Tea")
		return fmt.Errorf("unexpected model type: %T", finalModel)
	}

	actionName := "none"
	switch final.UserAction {
	case model.ActionAccept:
		actionName = "accept"
	case model.ActionEdit:
		actionName = "edit"
	case model.ActionRefresh:
		actionName = "refresh"
	case model.ActionCancel:
		actionName = "cancel"
	}

	t.log.Info().
		Str("flow", "commit").
		Str("action", actionName).
		Msg("TUI exited")

	if final.Err != nil {
		fmt.Println(errors.Format(final.Err))
		return nil
	}

	if final.UserAction == model.ActionAccept {
		t.log.Debug().Msg("Committing changes after TUI exit")
		output, err := t.git.Commit(t.ctx, final.Content)
		if err != nil {
			t.log.Error().Err(err).Msg("Failed to commit")
			return err
		}
		t.log.Info().Msg("Changes committed successfully")
		fmt.Println("🚀 Changes committed:")
		fmt.Println(output)
	}

	return nil
}
