package gpt

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

type stubLogger struct {
	log zerolog.Logger
}

func newStubLogger() stubLogger {
	return stubLogger{
		log: zerolog.Nop(),
	}
}

func (s stubLogger) Debug() *zerolog.Event        { return s.log.Debug() }
func (s stubLogger) Info() *zerolog.Event         { return s.log.Info() }
func (s stubLogger) Warn() *zerolog.Event         { return s.log.Warn() }
func (s stubLogger) Error() *zerolog.Event        { return s.log.Error() }
func (s stubLogger) GetBuffer() *logger.LogBuffer { return nil }

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_NewGPTClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	cfg.Logging.Level = "error"
	log := newStubLogger()

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "missing API token",
			token:       "",
			expectError: true,
		},
		{
			name:        "valid API token",
			token:       "valid-token",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldToken := os.Getenv("OPENAI_API_KEY")
			defer os.Setenv("OPENAI_API_KEY", oldToken)

			os.Setenv("OPENAI_API_KEY", tt.token)

			clientInstance, err := NewGPTClient(cfg, log)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, clientInstance)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, clientInstance)
				assert.IsType(t, &client{}, clientInstance)
			}
		})
	}
}

func Test_FetchCommitMessage(t *testing.T) {
	cfg := &config.Config{}
	cfg.Model.Name = "gpt-4"
	cfg.Model.MaxTokens = 500
	cfg.Model.Temperature = 0.7
	cfg.API.RetryCount = 0
	cfg.API.Timeout = 60000000000
	cfg.Logging.Level = "error"
	log := newStubLogger()

	tests := []struct {
		name        string
		diff        string
		before      func(*MockAPI)
		expected    string
		expectError bool
		errorType   error
	}{
		{
			name: "successful commit message generation",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{Message: openai.ChatCompletionMessage{Content: `{"type":"feat","scope":"api","description":"add endpoint","body":""}`}},
						},
					},
					nil,
				)
			},
			expected:    "feat(api): add endpoint",
			expectError: false,
		},
		{
			name: "API error",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{},
					errors.New("API error"),
				).AnyTimes()
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "empty response",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{}},
					nil,
				)
			},
			expected:    "",
			expectError: true,
			errorType:   errors.ErrNoResponse,
		},
		{
			name: "empty content in response",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{Message: openai.ChatCompletionMessage{Content: ""}},
						},
					},
					nil,
				)
			},
			expected:    "",
			expectError: true,
			errorType:   errors.ErrNoResponse,
		},
		{
			name: "invalid JSON in response",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{Message: openai.ChatCompletionMessage{Content: "invalid json"}},
						},
					},
					nil,
				)
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAPI := NewMockAPI(ctrl)
			tt.before(mockAPI)

			c := &client{
				cfg: cfg,
				api: mockAPI,
				log: log,
			}

			result, err := c.FetchCommitMessage(context.Background(), tt.diff)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_FetchChangelog(t *testing.T) {
	cfg := &config.Config{}
	cfg.Model.Name = "gpt-4"
	cfg.Model.MaxTokens = 500
	cfg.Model.Temperature = 0.7
	cfg.API.RetryCount = 0
	cfg.API.Timeout = 60000000000
	cfg.Logging.Level = "error"
	log := newStubLogger()

	tests := []struct {
		name        string
		commits     string
		before      func(*MockAPI)
		expected    string
		expectError bool
		errorType   error
	}{
		{
			name:    "successful changelog generation",
			commits: "abc123 feat: add feature\ndef456 fix: fix bug",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{Message: openai.ChatCompletionMessage{Content: "# CHANGELOG\n\n## Features\n- add feature"}},
						},
					},
					nil,
				)
			},
			expected:    "# CHANGELOG\n\n## Features\n- add feature",
			expectError: false,
		},
		{
			name:    "API error",
			commits: "abc123 feat: add feature",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{},
					errors.New("API error"),
				).AnyTimes()
			},
			expected:    "",
			expectError: true,
		},
		{
			name:    "empty response",
			commits: "abc123 feat: add feature",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{}},
					nil,
				)
			},
			expected:    "",
			expectError: true,
			errorType:   errors.ErrNoResponse,
		},
		{
			name:    "empty content in response",
			commits: "abc123 feat: add feature",
			before: func(mockAPI *MockAPI) {
				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{Message: openai.ChatCompletionMessage{Content: ""}},
						},
					},
					nil,
				)
			},
			expected:    "",
			expectError: true,
			errorType:   errors.ErrNoResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAPI := NewMockAPI(ctrl)
			tt.before(mockAPI)

			c := &client{
				cfg: cfg,
				api: mockAPI,
				log: log,
			}

			result, err := c.FetchChangelog(context.Background(), tt.commits)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_parseCommitMessageResponse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "simple commit message without scope",
			input:       `{"type":"feat","scope":"","description":"add new feature","body":""}`,
			expected:    "feat: add new feature",
			expectError: false,
		},
		{
			name:        "commit message with scope",
			input:       `{"type":"fix","scope":"api","description":"fix bug in endpoint","body":""}`,
			expected:    "fix(api): fix bug in endpoint",
			expectError: false,
		},
		{
			name:        "commit message with body",
			input:       `{"type":"feat","scope":"ui","description":"add button","body":"This adds a new button to the UI"}`,
			expected:    "feat(ui): add button\n\nThis adds a new button to the UI",
			expectError: false,
		},
		{
			name:        "commit message wrapped in json code block",
			input:       "```json\n{\"type\":\"chore\",\"scope\":\"\",\"description\":\"update dependencies\",\"body\":\"\"}\n```",
			expected:    "chore: update dependencies",
			expectError: false,
		},
		{
			name:        "invalid json",
			input:       `{"type":"feat"`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   \n  \t  ",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCommitMessageResponse(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_shouldRetry(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "rate limit error (429) - should retry",
			err:      &openai.APIError{HTTPStatusCode: 429},
			expected: true,
		},
		{
			name:     "internal server error (500) - should retry",
			err:      &openai.APIError{HTTPStatusCode: 500},
			expected: true,
		},
		{
			name:     "bad gateway (502) - should retry",
			err:      &openai.APIError{HTTPStatusCode: 502},
			expected: true,
		},
		{
			name:     "service unavailable (503) - should retry",
			err:      &openai.APIError{HTTPStatusCode: 503},
			expected: true,
		},
		{
			name:     "gateway timeout (504) - should retry",
			err:      &openai.APIError{HTTPStatusCode: 504},
			expected: true,
		},
		{
			name:     "unauthorized (401) - should not retry",
			err:      &openai.APIError{HTTPStatusCode: 401},
			expected: false,
		},
		{
			name:     "forbidden (403) - should not retry",
			err:      &openai.APIError{HTTPStatusCode: 403},
			expected: false,
		},
		{
			name:     "bad request (400) - should not retry",
			err:      &openai.APIError{HTTPStatusCode: 400},
			expected: false,
		},
		{
			name:     "non-API error - should retry",
			err:      errors.New("network error"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldRetry(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_validate(t *testing.T) {
	tests := []struct {
		name   string
		schema struct {
			Type        string `json:"type"`
			Scope       string `json:"scope"`
			Description string `json:"description"`
			Body        string `json:"body"`
		}
		expectError bool
		errorType   error
	}{
		{
			name: "valid schema with all fields",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{
				Type:        "feat",
				Scope:       "api",
				Description: "add endpoint",
				Body:        "detailed description",
			},
			expectError: false,
		},
		{
			name: "valid schema without scope and body",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{
				Type:        "fix",
				Scope:       "",
				Description: "fix bug",
				Body:        "",
			},
			expectError: false,
		},
		{
			name: "missing type",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{
				Type:        "",
				Scope:       "api",
				Description: "add endpoint",
				Body:        "",
			},
			expectError: true,
			errorType:   errors.ErrMissingCommitType,
		},
		{
			name: "missing description",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{
				Type:        "feat",
				Scope:       "api",
				Description: "",
				Body:        "",
			},
			expectError: true,
			errorType:   errors.ErrMissingCommitDesc,
		},
		{
			name: "invalid commit type",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{
				Type:        "invalid",
				Scope:       "api",
				Description: "add endpoint",
				Body:        "",
			},
			expectError: true,
			errorType:   errors.ErrInvalidCommitType,
		},
		{
			name: "all valid commit types - feat",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "feat", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - fix",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "fix", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - build",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "build", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - chore",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "chore", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - ci",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "ci", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - docs",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "docs", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - style",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "style", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - refactor",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "refactor", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - perf",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "perf", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - test",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "test", Description: "test"},
			expectError: false,
		},
		{
			name: "all valid commit types - revert",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "revert", Description: "test"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.schema)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
