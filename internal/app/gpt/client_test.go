package gpt

import (
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_NewHTTPClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.DefaultConfig()
	mockLogger := logger.NewMockLogger(ctrl)
	nopLogger := zerolog.Nop()
	mockEvent := nopLogger.Debug()

	tests := []struct {
		name     string
		before   func()
		token    string
		logLevel string
		error    bool
	}{
		{
			name: "Success with debug logging",
			before: func() {
				cfg.Logging.Level = "debug"
				mockLogger.EXPECT().Debug().Return(mockEvent).AnyTimes()
			},
			token:    "test-token",
			logLevel: "debug",
			error:    false,
		},
		{
			name: "Success with info logging",
			before: func() {
				cfg.Logging.Level = "info"
				mockLogger.EXPECT().Info().Return(mockEvent).AnyTimes()
			},
			token:    "test-token",
			logLevel: "info",
			error:    false,
		},
		{
			name: "wNo API token",
			before: func() {
				cfg.Logging.Level = "info"
				mockLogger.EXPECT().Error().Return(mockEvent).AnyTimes()
			},
			token:    "",
			logLevel: "info",
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			originalToken := os.Getenv("OPENAI_API_KEY")
			defer os.Setenv("OPENAI_API_KEY", originalToken)
			os.Setenv("OPENAI_API_KEY", tt.token)

			httpClient, err := NewHTTPClient(cfg, mockLogger)
			if tt.error {
				assert.Error(t, err)
				assert.Nil(t, httpClient)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, httpClient)

				restyClient, ok := httpClient.(*resty.Client)
				assert.True(t, ok)
				assert.Equal(t, BaseURL, restyClient.BaseURL)
			}
		})
	}
}
