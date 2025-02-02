package gpt

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	"cmt/internal/config"
)

// BaseURL the base URL for the OpenAI API
const BaseURL = "https://api.openai.com/v1"

// HTTPClient represents the HTTP client interface
type HTTPClient interface {
	R() *resty.Request
	SetBaseURL(url string) *resty.Client
	SetHeader(header, value string) *resty.Client
	SetRetryCount(count int) *resty.Client
}

// NewClient creates a new HTTP client
func NewClient() (HTTPClient, error) {
	token, err := config.GetAPIToken()
	if err != nil {
		log.Printf("🔐 Error: %s\n", err)
		return nil, err
	}

	return resty.New().
		SetBaseURL(BaseURL).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3), nil
}
