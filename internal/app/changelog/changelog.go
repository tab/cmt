package changelog

import (
	"context"
	"fmt"
	"time"

	"cmt/internal/app/cli/workflow"
	"cmt/internal/app/errors"
	"cmt/internal/config/logger"
)

// Generator handles changelog generation with dependencies injected via FX
type Generator struct {
	ctx      context.Context
	workflow workflow.Service
	log      logger.Logger
}

// NewGenerator creates a new changelog generator
func NewGenerator(
	ctx context.Context,
	workflow workflow.Service,
	log logger.Logger,
) *Generator {
	return &Generator{
		ctx:      ctx,
		workflow: workflow,
		log:      log,
	}
}

// Generate generates and outputs the changelog
func (g *Generator) Generate(between string) error {
	if between != "" {
		g.log.Info().Str("range", between).Msg("Generating changelog")
	} else {
		g.log.Info().Msg("Generating changelog from all commits")
	}

	spinner := startSpinner("Loading...")

	result, err := g.workflow.GenerateChangelog(g.ctx, between)

	stopSpinner(spinner)

	if err != nil {
		g.log.Error().Err(err).Msg("Failed to generate changelog")
		fmt.Println(errors.Format(err))
		return err
	}

	fmt.Println(result.Content)
	g.log.Info().Msg("Changelog generated successfully")
	return nil
}

// startSpinner starts a console spinner with the given message
func startSpinner(message string) chan bool {
	done := make(chan bool)

	go func() {
		spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				// Clear the spinner line
				fmt.Print("\r\033[K")
				return
			case <-ticker.C:
				fmt.Printf("\r%s %s", spinChars[i], message)
				i = (i + 1) % len(spinChars)
			}
		}
	}()

	return done
}

// stopSpinner stops the spinner
func stopSpinner(done chan bool) {
	done <- true
	close(done)
	time.Sleep(10 * time.Millisecond)
}
