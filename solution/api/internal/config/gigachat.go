package config

import (
	"errors"
	"os"
)

var (
	errGigaChatClientSecretNotSet = errors.New("giga chat client secret not set")
)

const (
	gigaChatClientSecretEnvName = "GIGA_CHAT_CLIENT_SECRET"
)

// GigaChatConfig defines an interface for giga chat config.
type GigaChatConfig interface {
	ClientSecret() string
}

type gigaChatConfig struct {
	clientSecret string
}

// NewGigaChatConfig initializes a new GigaChat configuration from environment variables.
func NewGigaChatConfig() (GigaChatConfig, error) {
	gigaChatClientSecret := os.Getenv(gigaChatClientSecretEnvName)
	if len(gigaChatClientSecret) == 0 {
		return nil, errGigaChatClientSecretNotSet
	}

	return &gigaChatConfig{
		clientSecret: gigaChatClientSecret,
	}, nil
}

// ClientSecret returns the giga chat api client secret.
func (cfg *gigaChatConfig) ClientSecret() string {
	return cfg.clientSecret
}
