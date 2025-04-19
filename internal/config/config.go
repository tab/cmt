package config

import (
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
	DefaultEditor      = "vim"

	Version = "0.5.0"
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

	cfg.Editor = DefaultEditor

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

	token, err := GetAPIToken()
	if err != nil {
		return nil, err
	}

	_ = token

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
