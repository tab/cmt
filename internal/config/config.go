package config

import (
	"fmt"
	"os"
	"time"

	"cmt/internal/errors"
)

const (
	Timeout = 60 * time.Second
)

func GetAPIToken() (string, error) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return "", fmt.Errorf("%w: OPENAI_API_KEY environment variable not set", errors.ErrAPITokenNotSet)
	}

	return token, nil
}
