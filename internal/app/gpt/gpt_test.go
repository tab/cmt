package gpt

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_NewGPTClient(t *testing.T) {
	cfg := &config.Config{}
	cfg.Logging.Level = "error"
	nopLogger := zerolog.Nop()

	tests := []struct {
		name        string
		token       string
		before      func(*logger.MockLogger)
		expectError bool
	}{
		{
			name:  "Failure with missing token",
			token: "",
			before: func(mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)
			},
			expectError: true,
		},
		{
			name:        "Success with valid token",
			token:       "valid-token",
			before:      func(mockLogger *logger.MockLogger) {},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLogger := logger.NewMockLogger(ctrl)
			tt.before(mockLogger)

			t.Setenv("OPENAI_API_KEY", tt.token)

			clientInstance, err := NewGPTClient(cfg, mockLogger)
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
	nopLogger := zerolog.Nop()

	tests := []struct {
		name        string
		diff        string
		before      func(*MockAPI, *logger.MockLogger)
		expected    string
		expectError bool
		errorType   error
	}{
		{
			name: "Success with valid response",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(2)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()

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
			name: "Failure when a p i error",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).AnyTimes()
				mockLogger.EXPECT().Warn().Return(nopLogger.Warn()).AnyTimes()

				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{},
					errors.New("API error"),
				).AnyTimes()
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "Failure with empty response",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)

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
			name: "Failure with empty content",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)

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
			name: "Failure with invalid j s o n",
			diff: "diff --git a/file.go",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)

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
			mockLogger := logger.NewMockLogger(ctrl)
			tt.before(mockAPI, mockLogger)

			c := &client{
				cfg: cfg,
				api: mockAPI,
				log: mockLogger,
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
	nopLogger := zerolog.Nop()

	tests := []struct {
		name        string
		commits     string
		before      func(*MockAPI, *logger.MockLogger)
		expected    string
		expectError bool
		errorType   error
	}{
		{
			name:    "Success with valid response",
			commits: "abc123 feat: add feature\ndef456 fix: fix bug",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(2)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()

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
			name:    "Failure when a p i error",
			commits: "abc123 feat: add feature",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).AnyTimes()
				mockLogger.EXPECT().Warn().Return(nopLogger.Warn()).AnyTimes()

				mockAPI.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
					openai.ChatCompletionResponse{},
					errors.New("API error"),
				).AnyTimes()
			},
			expected:    "",
			expectError: true,
		},
		{
			name:    "Failure with empty response",
			commits: "abc123 feat: add feature",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)

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
			name:    "Failure with empty content",
			commits: "abc123 feat: add feature",
			before: func(mockAPI *MockAPI, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info().Return(nopLogger.Info()).Times(1)
				mockLogger.EXPECT().Debug().Return(nopLogger.Debug()).AnyTimes()
				mockLogger.EXPECT().Error().Return(nopLogger.Error()).Times(1)

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
			mockLogger := logger.NewMockLogger(ctrl)
			tt.before(mockAPI, mockLogger)

			c := &client{
				cfg: cfg,
				api: mockAPI,
				log: mockLogger,
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

func Test_ParseCommitMessageResponse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Success without scope",
			input:       `{"type":"feat","scope":"","description":"add new feature","body":""}`,
			expected:    "feat: add new feature",
			expectError: false,
		},
		{
			name:        "Success with scope",
			input:       `{"type":"fix","scope":"api","description":"fix bug in endpoint","body":""}`,
			expected:    "fix(api): fix bug in endpoint",
			expectError: false,
		},
		{
			name:        "Success with body",
			input:       `{"type":"feat","scope":"ui","description":"add button","body":"This adds a new button to the UI"}`,
			expected:    "feat(ui): add button\n\nThis adds a new button to the UI",
			expectError: false,
		},
		{
			name:        "Success with j s o n code block",
			input:       "```json\n{\"type\":\"chore\",\"scope\":\"\",\"description\":\"update dependencies\",\"body\":\"\"}\n```",
			expected:    "chore: update dependencies",
			expectError: false,
		},
		{
			name:        "Failure with invalid j s o n",
			input:       `{"type":"feat"`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Failure with empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Failure with whitespace only",
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

func Test_ShouldRetry(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Success when rate limited",
			err:      &openai.APIError{HTTPStatusCode: 429},
			expected: true,
		},
		{
			name:     "Success when internal error",
			err:      &openai.APIError{HTTPStatusCode: 500},
			expected: true,
		},
		{
			name:     "Success when bad gateway",
			err:      &openai.APIError{HTTPStatusCode: 502},
			expected: true,
		},
		{
			name:     "Success when service unavailable",
			err:      &openai.APIError{HTTPStatusCode: 503},
			expected: true,
		},
		{
			name:     "Success when gateway timeout",
			err:      &openai.APIError{HTTPStatusCode: 504},
			expected: true,
		},
		{
			name:     "Success without retry on unauthorized",
			err:      &openai.APIError{HTTPStatusCode: 401},
			expected: false,
		},
		{
			name:     "Success without retry on forbidden",
			err:      &openai.APIError{HTTPStatusCode: 403},
			expected: false,
		},
		{
			name:     "Success without retry on bad request",
			err:      &openai.APIError{HTTPStatusCode: 400},
			expected: false,
		},
		{
			name:     "Success when non a p i error",
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

func Test_Validate(t *testing.T) {
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
			name: "Success with all fields",
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
			name: "Success without scope and body",
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
			name: "Failure with missing type",
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
			name: "Failure with missing description",
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
			name: "Failure with invalid commit type",
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
			name: "Success with type feat",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "feat", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type fix",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "fix", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type build",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "build", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type chore",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "chore", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type c i",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "ci", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type docs",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "docs", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type style",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "style", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type refactor",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "refactor", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type perf",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "perf", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type test",
			schema: struct {
				Type        string `json:"type"`
				Scope       string `json:"scope"`
				Description string `json:"description"`
				Body        string `json:"body"`
			}{Type: "test", Description: "test"},
			expectError: false,
		},
		{
			name: "Success with type revert",
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
