package config

import (
	"errors"
	"os"
	"strconv"
)

var (
	errAdModerationEnabledNotSet = errors.New("ad moderation enabled not set")
)

const (
	adModerationEnabledEnvName = "AD_MODERATION_ENABLED"
)

// ModerationConfig defines an interface for moderation config.
type ModerationConfig interface {
	AdModerationEnabled() bool
}

type moderationConfig struct {
	adModerationEnabled bool
}

// NewModerationConfig initializes a new Moderation configuration from environment variables.
func NewModerationConfig() (ModerationConfig, error) {
	adModerationEnabled, err := strconv.ParseBool(os.Getenv(adModerationEnabledEnvName))
	if err != nil {
		return nil, errAdModerationEnabledNotSet
	}

	return &moderationConfig{
		adModerationEnabled: adModerationEnabled,
	}, nil
}

// AdModerationEnabled returns .
func (cfg *moderationConfig) AdModerationEnabled() bool {
	return cfg.adModerationEnabled
}
