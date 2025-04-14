package config

import (
	"errors"
	"os"
)

var (
	errAdvertisingURLNotSet = errors.New("advertising url not set")
)

const (
	advertisingURLEnvName = "ADVERTISING_URL"
)

// AdvertisingConfig defines an interface for advertising platform configuration.
type AdvertisingConfig interface {
	URL() string
}

type advertisingConfig struct {
	url string
}

// NewAdvertisingConfig loads advertising platform configuration from environment variables.
func NewAdvertisingConfig() (AdvertisingConfig, error) {
	url := os.Getenv(advertisingURLEnvName)
	if len(url) == 0 {
		return nil, errAdvertisingURLNotSet
	}

	return &advertisingConfig{
		url: url,
	}, nil
}

func (cfg *advertisingConfig) URL() string {
	return cfg.url
}
