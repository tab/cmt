package gpt

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"cmt/internal/config"
	"cmt/internal/config/logger"
)

// BaseURL is the base URL for the OpenAI API
const BaseURL = "https://api.openai.com/v1"

// HTTPClient represents the HTTP client interface
type HTTPClient interface {
	R() *resty.Request
	SetBaseURL(url string) *resty.Client
	SetHeader(header, value string) *resty.Client
	SetRetryCount(count int) *resty.Client
}

// NewHTTPClient creates a new HTTP client for the OpenAI API
func NewHTTPClient(
	cfg *config.Config,
	log logger.Logger,
) (HTTPClient, error) {
	token, err := config.GetAPIToken()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get API token")
		return nil, err
	}

	httpClient := resty.New().
		SetBaseURL(BaseURL).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(cfg.API.RetryCount)

	log.Debug().
		Int("retry_count", cfg.API.RetryCount).
		Str("base_url", BaseURL).
		Msg("HTTP client created")

	return httpClient, nil
}
