package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cmt/internal/errors"
)

const (
	ModelName   = "gpt-4o-mini"
	MaxTokens   = 500
	Temperature = 0.7
)

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTRequest struct {
	Model       string       `json:"model"`
	Messages    []GPTMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
}

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type GPTModelClient interface {
	FetchCommitMessage(ctx context.Context, diff string) (string, error)
	FetchChangelog(ctx context.Context, commits string) (string, error)
}

type GPTModel struct {
	Client HTTPClient
}

func NewGPTModel() *GPTModel {
	return &GPTModel{}
}

func (g *GPTModel) FetchCommitMessage(ctx context.Context, diff string) (string, error) {
	const systemPrompt = `You are an experienced Software Engineer tasked with generating a Conventional Commit message based on a provided git diff.
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
}
`

	messages := []GPTMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: diff},
	}

	content, err := g.fetch(ctx, messages)
	if err != nil {
		return "", err
	}

	message, err := parseCommitMessageResponse(content)
	if err != nil {
		return "", err
	}

	return message, nil
}

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
		return "", fmt.Errorf("Failed to parse JSON: %w", errors.ErrFailedToParseJSON)
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

func (g *GPTModel) FetchChangelog(ctx context.Context, commits string) (string, error) {
	const systemPrompt = `You are an experienced Software Engineer tasked with generating a concise and clear CHANGELOG for a set of commits in Markdown format.
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
### Bug Fixes
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
- **BREAKING CHANGE:** ...
`

	messages := []GPTMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: commits},
	}

	content, err := g.fetch(ctx, messages)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return content, nil
}

func (g *GPTModel) fetch(ctx context.Context, messages []GPTMessage) (string, error) {
	if g.Client == nil {
		var err error
		g.Client, err = NewClient()
		if err != nil {
			return "", err
		}
	}

	requestBody := GPTRequest{
		Model:       ModelName,
		Messages:    messages,
		MaxTokens:   MaxTokens,
		Temperature: Temperature,
	}

	resp, err := g.Client.R().
		SetContext(ctx).
		SetBody(requestBody).
		Post("/chat/completions")

	if err != nil {
		return "", fmt.Errorf("API request error: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("API error: %s", resp.String())
	}

	var result GPTResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", errors.ErrNoResponse
	}

	content := strings.TrimSpace(result.Choices[0].Message.Content)
	return content, nil
}
