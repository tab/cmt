package gpt

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_FetchCommitMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := config.DefaultConfig()
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	httpClient := resty.New()
	httpClient.SetTransport(httpmock.DefaultTransport)
	httpClient.SetBaseURL(BaseURL)
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetRetryCount(3)

	gptClient := NewGPTClient(cfg, httpClient, mockLogger)

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
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()

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
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()

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
			name: "Empty message in choices",
			before: func() {
				httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
					httpmock.NewStringResponder(200, `{
						"choices": [{
							"message": {}
						}]
					}`))
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
			tt.before()

			httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", args.token))
			content, err := gptClient.FetchCommitMessage(context.Background(), args.diff)

			if tt.expected.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.content, content)
			}
		})
	}
}

func Test_FetchChangelog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := config.DefaultConfig()
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	httpClient := resty.New()
	httpClient.SetTransport(httpmock.DefaultTransport)
	httpClient.SetBaseURL(BaseURL)
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetRetryCount(3)

	gptClient := NewGPTClient(cfg, httpClient, mockLogger)

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
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()

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
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()

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
			tt.before()

			httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", args.token))
			content, err := gptClient.FetchChangelog(context.Background(), args.commits)

			if tt.expected.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.content, content)
			}
		})
	}
}

func Test_parseCommitMessageResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected string
		error    bool
	}{
		{
			name:     "Valid JSON string",
			response: "```json\n{\"type\": \"feat\", \"scope\": \"core\", \"description\": \"Add new feature\", \"body\": \"Detailed explanation\"}\n```",
			expected: "feat(core): Add new feature\n\nDetailed explanation",
			error:    false,
		},
		{
			name:     "Valid JSON without code block",
			response: "{\"type\": \"fix\", \"scope\": \"api\", \"description\": \"Fix bug\", \"body\": \"Bug details\"}",
			expected: "fix(api): Fix bug\n\nBug details",
			error:    false,
		},
		{
			name:     "Valid JSON with empty scope",
			response: "{\"type\": \"chore\", \"scope\": \"\", \"description\": \"Update dependencies\", \"body\": \"Updated all dependencies\"}",
			expected: "chore: Update dependencies\n\nUpdated all dependencies",
			error:    false,
		},
		{
			name:     "Valid JSON with empty body",
			response: "{\"type\": \"docs\", \"scope\": \"\", \"description\": \"Update README\", \"body\": \"\"}",
			expected: "docs: Update README",
		},
		{
			name:     "Invalid JSON",
			response: `invalid json`,
			expected: "",
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCommitMessageResponse(tt.response)

			if tt.error {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
