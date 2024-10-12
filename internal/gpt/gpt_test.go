package gpt

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"cc-gpt/internal/errors"
)

func Test_BuildPrompt(t *testing.T) {
	type args struct {
		diff string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "With diff content",
			args: args{
				diff: "diff --git a/file.go b/file.go\nindex 83db48f..f5f6f70 100644\n--- a/file.go\n+++ b/file.go\n@@ -1,4 +1,4 @@\n-package main\n+package main_test",
			},
			expected: `Generate a conventional commit message in JSON format with the following structure: { "type": "<type>", "scope": "<scope>", "description": "<description>", "body": "<body>" }. Use git diff: diff --git a/file.go b/file.go
index 83db48f..f5f6f70 100644
--- a/file.go
+++ b/file.go
@@ -1,4 +1,4 @@
-package main
+package main_test`,
		},
		{
			name: "Empty diff",
			args: args{
				diff: "",
			},
			expected: `Generate a conventional commit message in JSON format with the following structure: { "type": "<type>", "scope": "<scope>", "description": "<description>", "body": "<body>" }. Use git diff: `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPrompt(tt.args.diff)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_ParseContent(t *testing.T) {
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
			input: `{"type": "feat", "scope": "test", "description": "Add new feature", "body": "Details about the feature"}`,
			expected: result{
				content: "feat(test): Add new feature\n\nDetails about the feature",
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
			result, err := parseContent(tt.input)

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

func Test_Fetch(t *testing.T) {
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
			result, err := model.Fetch(context.Background(), args.diff)

			if tt.expected.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.content, result)
			}
		})
	}
}
