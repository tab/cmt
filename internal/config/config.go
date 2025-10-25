package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"cmt/internal/app/errors"
)

const (
	DefaultTimeout     = 60 * time.Second
	DefaultModelName   = "gpt-4.1-nano"
	DefaultMaxTokens   = 500
	DefaultTemperature = 0.7
	DefaultRetryCount  = 3
	DefaultLogLevel    = "info"
	DefaultLogFormat   = "console"

	AppName        = "cmt"
	AppDescription = "command line utility to generate conversational commits using OpenAI's GPT models"

	Version = "0.7.0"
)

// Config represents the application configuration
type Config struct {
	Model struct {
		Name        string  `yaml:"name"`
		MaxTokens   int     `yaml:"max_tokens"`
		Temperature float64 `yaml:"temperature"`
	} `yaml:"model"`
	API struct {
		RetryCount int           `yaml:"retry_count"`
		Timeout    time.Duration `yaml:"timeout"`
	} `yaml:"api"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	}
	Editor string `yaml:"editor"`
}

// Option allows for functional options pattern
type Option func(*Config)

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cfg := &Config{}

	cfg.Model.Name = DefaultModelName
	cfg.Model.MaxTokens = DefaultMaxTokens
	cfg.Model.Temperature = DefaultTemperature

	cfg.API.RetryCount = DefaultRetryCount
	cfg.API.Timeout = DefaultTimeout

	cfg.Logging.Level = DefaultLogLevel
	cfg.Logging.Format = DefaultLogFormat

	return cfg
}

// Load loads the configuration from file
func Load() (*Config, error) {
	cfg := DefaultConfig()

	v := viper.New()
	v.SetConfigName("cmt")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.ErrFailedToReadConfig
		}
	} else {
		if err = v.Unmarshal(cfg); err != nil {
			return nil, errors.ErrFailedToParseConfig
		}
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetAPIToken returns the OpenAI API token
func GetAPIToken() (string, error) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return "", errors.ErrAPITokenNotSet
	}

	return token, nil
}

// validate validates the configuration values
func (c *Config) validate() error {
	if c.Model.Temperature < 0 || c.Model.Temperature > 2 {
		return fmt.Errorf("%w: must be between 0 and 2, got %.2f", errors.ErrInvalidTemperature, c.Model.Temperature)
	}
	if c.Model.MaxTokens <= 0 {
		return fmt.Errorf("%w: must be positive, got %d", errors.ErrInvalidMaxTokens, c.Model.MaxTokens)
	}
	if c.API.Timeout <= 0 {
		return fmt.Errorf("%w: must be positive, got %v", errors.ErrInvalidTimeout, c.API.Timeout)
	}
	if c.API.RetryCount < 0 {
		return fmt.Errorf("%w: must be non-negative, got %d", errors.ErrInvalidRetryCount, c.API.RetryCount)
	}
	return nil
}
