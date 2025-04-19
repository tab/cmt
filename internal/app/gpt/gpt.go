package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

const (
	CommitSystemPromt = `You are an experienced Software Engineer tasked with generating a Conventional Commit message based on a provided git diff.
Follow these guidelines:
– Use the Conventional Commits format for the commit message
– Do not elaborate unnecessarily. Focus on the core details
– Ensure proper formatting for easy readability

Follow the format below:
{
  "type": "feat, fix, build, chore, ci, docs, style, refactor, perf, test, revert",
  "scope": "scope of the change (use one word)",
  "description": "a brief description of what was changed",
  "body": "an optional longer explanation of the change",
  "footer": "any additional information, like breaking changes or issue links"
}`

	ChangelogSystemPrompt = `You are an experienced Software Engineer tasked with generating a concise and clear CHANGELOG for a set of commits in Markdown format.
Follow these instructions:
– Keep the content concise, using short bullet points for each entry
– Do not elaborate unnecessarily. Focus on the core details
– Ensure proper formatting for easy readability
– Do not include empty sections

Follow the format below:

# CHANGELOG

## [X.Y.Z]

### Features
- **feat:** ...
### Fixes
- **fix:** ...
### Performance
- **perf:** ...
### Refactor
- **refactor:** ...
### Documentation
- **docs:** ...
### Style
- **style:** ...
### Tests
- **test:** ...
### Build
- **build:** ...
### CI
- **ci:** ...
### Chore
- **chore:** ...
### Reverts
- **revert:** ...
### Breaking Changes
- **BREAKING CHANGE:** ...`
)

// Client represents the GPT model client interface
type Client interface {
	FetchCommitMessage(ctx context.Context, diff string) (string, error)
	FetchChangelog(ctx context.Context, commits string) (string, error)
}

// Message represents a GPT message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents a GPT request
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Response represents a GPT response
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// client implements the Client interface
type client struct {
	cfg        *config.Config
	httpClient HTTPClient
	log        logger.Logger
}

// NewGPTClient creates a new GPT model client
func NewGPTClient(
	cfg *config.Config,
	httpClient HTTPClient,
	log logger.Logger,
) Client {
	return &client{
		cfg:        cfg,
		httpClient: httpClient,
		log:        log,
	}
}

// FetchCommitMessage generates a commit message from a git diff
func (g *client) FetchCommitMessage(ctx context.Context, diff string) (string, error) {
	messages := []Message{
		{Role: "system", Content: CommitSystemPromt},
		{Role: "user", Content: diff},
	}

	g.log.Debug().Msg("Fetching commit message from GPT")
	content, err := g.fetch(ctx, messages)
	if err != nil {
		g.log.Debug().Err(err).Msg("Failed to fetch commit message")
		return "", err
	}

	message, err := parseCommitMessageResponse(content)
	if err != nil {
		g.log.Error().Err(err).Msg("Failed to parse commit message response")
		return "", err
	}

	return message, nil
}

// FetchChangelog generates a changelog from git commits
func (g *client) FetchChangelog(ctx context.Context, commits string) (string, error) {
	messages := []Message{
		{Role: "system", Content: ChangelogSystemPrompt},
		{Role: "user", Content: commits},
	}

	g.log.Debug().Msg("Fetching changelog from GPT")
	content, err := g.fetch(ctx, messages)
	if err != nil {
		g.log.Error().Err(err).Msg("Failed to fetch changelog")
		return "", err
	}

	return content, nil
}

// fetch sends a request to the GPT API and returns the response content
func (g *client) fetch(ctx context.Context, messages []Message) (string, error) {
	requestBody := Request{
		Model:       g.cfg.Model.Name,
		Messages:    messages,
		MaxTokens:   g.cfg.Model.MaxTokens,
		Temperature: g.cfg.Model.Temperature,
	}

	g.log.Debug().
		Str("model", g.cfg.Model.Name).
		Int("max_tokens", g.cfg.Model.MaxTokens).
		Float64("temperature", g.cfg.Model.Temperature).
		Msg("Sending GPT request")

	resp, err := g.httpClient.R().
		SetContext(ctx).
		SetBody(requestBody).
		Post("/chat/completions")

	if err != nil {
		return "", fmt.Errorf("API request error: %w", err)
	}

	if resp.IsError() {
		g.log.Debug().
			Int("status_code", resp.StatusCode()).
			Str("response", resp.String()).
			Msg("API returned error")
		return "", fmt.Errorf("API error: %s (status code: %d)", resp.String(), resp.StatusCode())
	}

	var result Response
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", errors.ErrNoResponse
	}

	content := strings.TrimSpace(result.Choices[0].Message.Content)
	g.log.Debug().Msg("Successfully received GPT response")
	return content, nil
}

// parseCommitMessageResponse parses the GPT response into a conventional commit message
func parseCommitMessageResponse(text string) (string, error) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}

	var schema struct {
		Type        string `json:"type"`
		Scope       string `json:"scope"`
		Description string `json:"description"`
		Body        string `json:"body"`
	}

	if err := json.Unmarshal([]byte(text), &schema); err != nil {
		return "", fmt.Errorf("%w: %v", errors.ErrFailedToParseJSON, err)
	}

	scope := ""
	if schema.Scope != "" {
		scope = fmt.Sprintf("(%s)", schema.Scope)
	}

	header := fmt.Sprintf("%s%s: %s", schema.Type, scope, schema.Description)
	message := header
	if schema.Body != "" {
		message = fmt.Sprintf("%s\n\n%s", header, schema.Body)
	}

	return message, nil
}
