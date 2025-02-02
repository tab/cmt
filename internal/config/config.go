package config

import (
	"fmt"
	"os"
	"time"

	"cmt/internal/errors"
)

// Timeout the default timeout
const (
	Timeout = 60 * time.Second
)

// GetAPIToken returns the OpenAI API token
func GetAPIToken() (string, error) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return "", fmt.Errorf("%w: OPENAI_API_KEY environment variable not set", errors.ErrAPITokenNotSet)
	}

	return token, nil
}
