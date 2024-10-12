package gpt

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	"cmt/internal/config"
)

type HTTPClient interface {
	R() *resty.Request
	SetBaseURL(url string) *resty.Client
	SetHeader(header, value string) *resty.Client
	SetRetryCount(count int) *resty.Client
}

func NewClient() (HTTPClient, error) {
	token, err := config.GetAPIToken()
	if err != nil {
		log.Fatalf("🔐 Error: %s\n", err)
	}

	return resty.New().
		SetBaseURL(BASE_URL).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3), nil
}
