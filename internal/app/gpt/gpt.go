package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"

	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

const (
	commitSystemPromt = `You are an expert at writing Conventional Commit messages following the v1.0.0 specification.
Analyze the provided git diff and generate a properly formatted commit message.

RULES:
1. Type: Choose the most appropriate type based on the change:
   - feat: new feature or capability
   - fix: bug fix or correction
   - docs: documentation only changes
   - style: code style/formatting (no functional changes)
   - refactor: code restructuring (no functional changes)
   - perf: performance improvements
   - test: adding or updating tests
   - build: build system or dependencies
   - ci: CI/CD configuration changes
   - chore: other changes (tooling, configs)
   - revert: reverting a previous commit

2. Scope: A one word noun describing the codebase section or package (e.g., parser, api, auth)

3. Description:
   - Start with uppercase letter
   - Use imperative mood ("Add" not "Added" or "Adds")
   - No period at the end
   - Be concise but clear

4. Body: Provide additional context if the change is non-trivial.
   - Use to explain "what" and "why", not "how"
   - Wrap at 72 characters

EXAMPLES:
- feat(auth): Add OAuth2 login support
- docs(readme): Update installation instructions
- refactor(parser): Extract validation logic into separate function

Return ONLY valid JSON in this format:
{
  "type": "feat, fix, build, chore, ci, docs, style, refactor, perf, test, revert",
  "scope": "scope of the change (use one word)",
  "description": "a brief description of what was changed in imperative mood",
  "body": "optional detailed explanation"
}`

	changelogSystemPrompt = `You are an experienced Software Engineer tasked with generating a concise and clear CHANGELOG for a set of commits in Markdown format.
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
- **BREAKING CHANGES:** ...`
)

type commitType string

const (
	commitTypeFeat     commitType = "feat"
	commitTypeFix      commitType = "fix"
	commitTypeBuild    commitType = "build"
	commitTypeChore    commitType = "chore"
	commitTypeCI       commitType = "ci"
	commitTypeDocs     commitType = "docs"
	commitTypeStyle    commitType = "style"
	commitTypeRefactor commitType = "refactor"
	commitTypePerf     commitType = "perf"
	commitTypeTest     commitType = "test"
	commitTypeRevert   commitType = "revert"
)

var commitTypes = []commitType{
	commitTypeFeat,
	commitTypeFix,
	commitTypeBuild,
	commitTypeChore,
	commitTypeCI,
	commitTypeDocs,
	commitTypeStyle,
	commitTypeRefactor,
	commitTypePerf,
	commitTypeTest,
	commitTypeRevert,
}

// Client represents the GPT model client interface
type Client interface {
	FetchCommitMessage(ctx context.Context, diff string) (string, error)
	FetchChangelog(ctx context.Context, commits string) (string, error)
}

// API represents the OpenAI API client interface
type API interface {
	CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// client implements the Client interface
type client struct {
	cfg *config.Config
	api API
	log logger.Logger
}

// NewGPTClient creates a new GPT model client
func NewGPTClient(cfg *config.Config, log logger.Logger) (Client, error) {
	token, err := config.GetAPIToken()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get API token")
		return nil, err
	}

	clientCfg := openai.DefaultConfig(token)
	clientCfg.HTTPClient = &http.Client{
		Timeout: cfg.API.Timeout,
	}

	return &client{
		cfg: cfg,
		api: openai.NewClientWithConfig(clientCfg),
		log: log,
	}, nil
}

// FetchCommitMessage generates a commit message from a git diff
func (g *client) FetchCommitMessage(ctx context.Context, diff string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: commitSystemPromt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: diff,
		},
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
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: changelogSystemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: commits,
		},
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
func (g *client) fetch(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	g.log.Debug().
		Str("model", g.cfg.Model.Name).
		Int("max_tokens", g.cfg.Model.MaxTokens).
		Float64("temperature", g.cfg.Model.Temperature).
		Msg("Sending GPT request")

	if e := g.log.Debug(); e.Enabled() {
		payload, err := json.MarshalIndent(messages, "", "  ")
		if err != nil {
			e.Err(err).Msg("Failed to marshal request payload")
		} else {
			payloadStr := string(payload)
			if len(payloadStr) > 10000 {
				payloadStr = payloadStr[:10000] + "\n... (truncated)"
			}
			e.Msgf("Request payload:\n%s", payloadStr)
		}
	}

	var resp openai.ChatCompletionResponse
	var lastErr error

	for attempt := 0; attempt <= g.cfg.API.RetryCount; attempt++ {
		if attempt > 0 {
			backoffDuration := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			g.log.Debug().Int("attempt", attempt+1).Dur("backoff", backoffDuration).Msg("Retrying API request")

			select {
			case <-time.After(backoffDuration):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		var err error
		resp, err = g.api.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:               g.cfg.Model.Name,
				Messages:            messages,
				MaxCompletionTokens: g.cfg.Model.MaxTokens,
				Temperature:         float32(g.cfg.Model.Temperature),
			},
		)

		if err == nil {
			if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
				return "", errors.ErrNoResponse
			}

			content := strings.TrimSpace(resp.Choices[0].Message.Content)

			if e := g.log.Debug(); e.Enabled() {
				responseData, marshalErr := json.MarshalIndent(resp, "", "  ")
				if marshalErr != nil {
					e.Err(marshalErr).Msg("Failed to marshal API response")
				} else {
					responseStr := string(responseData)
					if len(responseStr) > 10000 {
						responseStr = responseStr[:10000] + "\n... (truncated)"
					}
					e.Msgf("API response:\n%s", responseStr)
				}
			}

			g.log.Debug().Msg("Successfully received GPT response")
			return content, nil
		}

		lastErr = err

		if shouldRetry(err) {
			if attempt < g.cfg.API.RetryCount {
				g.log.Warn().Err(err).Int("attempt", attempt+1).Msg("API request failed, will retry")
			}
			continue
		}

		g.log.Debug().Err(err).Msg("Non-retryable error, aborting")
		return "", fmt.Errorf("API request error: %w", err)
	}

	return "", fmt.Errorf("API request error after %d retries: %w", g.cfg.API.RetryCount, lastErr)
}

// shouldRetry determines if an error is retryable based on OpenAI API error codes
func shouldRetry(err error) bool {
	var apiErr *openai.APIError

	if errors.As(err, &apiErr) {
		switch apiErr.HTTPStatusCode {
		case 429:
			return true
		case 500, 502, 503, 504:
			return true
		default:
			return false
		}
	}

	return true
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

	if err := validate(schema); err != nil {
		return "", err
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

// validate validates the parsed commit message schema
func validate(schema struct {
	Type        string `json:"type"`
	Scope       string `json:"scope"`
	Description string `json:"description"`
	Body        string `json:"body"`
}) error {
	if schema.Type == "" {
		return errors.ErrMissingCommitType
	}

	if schema.Description == "" {
		return errors.ErrMissingCommitDesc
	}

	isValidType := false
	for _, validType := range commitTypes {
		if schema.Type == string(validType) {
			isValidType = true
			break
		}
	}

	if !isValidType {
		validTypes := make([]string, len(commitTypes))

		for i, t := range commitTypes {
			validTypes[i] = string(t)
		}

		return fmt.Errorf("%w: %q (must be one of: %v)", errors.ErrInvalidCommitType, schema.Type, validTypes)
	}

	return nil
}
