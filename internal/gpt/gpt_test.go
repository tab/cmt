package gpt

import (
	"cmt/internal/errors"
	"context"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_FetchCommitMessage(t *testing.T) {
	args := struct {
		diff  string
		token string
	}{
		diff:  "diff content",
		token: "test-token",
	}

	type result struct {
		content string
		error   bool
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `{
   				"choices": [{
   					"message": {
   						"content": "{\"type\": \"feat\", \"scope\": \"core\", \"description\": \"Add new feature\", \"body\": \"Detailed explanation.\"}"
   					}
   				}]
   			}`))
			},
			expected: result{
				content: "feat(core): Add new feature\n\nDetailed explanation.",
				error:   false,
			},
		},
		{
			name: "API error response",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(500, `{"error": "Internal Server Error"}`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
		{
			name: "Invalid JSON in response",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `invalid json`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
		{
			name: "No choices in response",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `{"choices": []}`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
		{
			name: "Empty content in message",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `{
						"choices": [{
							"message": {
								"content": ""
							}
						}]
					}`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			tt.before()

			client := resty.New()
			client.SetTransport(httpmock.DefaultTransport)
			client.SetBaseURL(BASE_URL)
			client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", args.token))
			client.SetHeader("Content-Type", "application/json")
			client.SetRetryCount(3)

			model := &GPTModel{Client: client}
			result, err := model.FetchCommitMessage(context.Background(), args.diff)

			if tt.expected.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.content, result)
			}
		})
	}
}

func Test_ParseCommitMessageResponse(t *testing.T) {
	type result struct {
		content string
		error   bool
	}

	tests := []struct {
		name     string
		input    string
		expected result
	}{
		{
			name:  "Valid JSON content",
			input: `{"type": "feat", "scope": "core", "description": "Add new feature", "body": "Details about the feature"}`,
			expected: result{
				content: "feat(core): Add new feature\n\nDetails about the feature",
				error:   false,
			},
		},
		{
			name:  "Invalid JSON content",
			input: `{"type": "feat", "scope": "test", "description"}`,
			expected: result{
				content: "",
				error:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCommitMessageResponse(tt.input)

			if tt.expected.error {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), errors.ErrFailedToParseJSON.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.content, result)
			}
		})
	}
}

func TestGPTModel_FetchChangelog(t *testing.T) {
	args := struct {
		commits string
		token   string
	}{
		commits: "abcd123 feat(jwt): Add new feature\n\nDetails about the feature",
		token:   "test-token",
	}

	type result struct {
		content string
		error   bool
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `{
            "choices": [{
              "message": {
                "content": "# CHANGELOG\n\n## Features\n\n- feat(jwt): Add new feature\n\nDetails about the feature"
              }
            }]
          }`))
			},
			expected: result{
				content: "# CHANGELOG\n\n## Features\n\n- feat(jwt): Add new feature\n\nDetails about the feature",
				error:   false,
			},
		},
		{
			name: "API error response",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(500, `{"error": "Internal Server Error"}`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
		{
			name: "Invalid JSON in response",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `invalid json`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
		{
			name: "No choices in response",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `{"choices": []}`))
			},
			expected: result{
				content: "",
				error:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			tt.before()

			client := resty.New()
			client.SetTransport(httpmock.DefaultTransport)
			client.SetBaseURL(BASE_URL)
			client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", args.token))
			client.SetHeader("Content-Type", "application/json")
			client.SetRetryCount(3)

			model := &GPTModel{Client: client}
			result, err := model.FetchChangelog(context.Background(), args.commits)

			if tt.expected.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.content, result)
			}
		})
	}
}
