package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cmt/internal/errors"
)

const (
	BASE_URL     = "https://api.openai.com/v1"
	MODEL_NAME   = "gpt-4o-mini"
	SYSTEM_PROMT = "You are a Software Engineer. Provide a Conventional Commit message for the git diff in JSON format."
	MAX_TOKENS   = 150
	TEMPERATURE  = 0.7
)

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTRequest struct {
	Model       string       `json:"model"`
	Messages    []GPTMessage `json:"messages"`
	MAX_TOKENS  int          `json:"max_tokens,omitempty"`
	TEMPERATURE float64      `json:"temperature,omitempty"`
}

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type GPTModel struct {
	Client HTTPClient
}

func (g *GPTModel) Fetch(ctx context.Context, diff string) (string, error) {
	if diff == "" {
		return "", errors.ErrNoGitChanges
	}

	if g.Client == nil {
		g.Client, _ = NewClient()
	}

	requestBody := GPTRequest{
		Model: MODEL_NAME,
		Messages: []GPTMessage{
			{
				Role:    "system",
				Content: SYSTEM_PROMT,
			},
			{
				Role:    "user",
				Content: buildPrompt(diff),
			},
		},
		MAX_TOKENS:  MAX_TOKENS,
		TEMPERATURE: TEMPERATURE,
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

	message, err := parseContent(strings.TrimSpace(result.Choices[0].Message.Content))
	if err != nil {
		return "", err
	}

	return message, nil
}

func buildPrompt(diff string) string {
	return fmt.Sprintf(`Generate a conventional commit message in JSON format with the following structure: { "type": "<type>", "scope": "<scope>", "description": "<description>", "body": "<body>" }. Use git diff: %s`, diff)
}

func parseContent(text string) (string, error) {
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
