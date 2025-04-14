package config

import (
	"errors"
	"os"
)

var (
	errBotTokenNotSet            = errors.New("bot token not set")
	errBotAuthBannerNotSet       = errors.New("bot auth banner not set")
	errBotClientBannerNotSet     = errors.New("bot client banner not set")
	errBotAdvertiserBannerNotSet = errors.New("bot advertiser banner not set")
)

const (
	botTokenEnvName            = "BOT_TOKEN"
	botAuthBannerEnvName       = "BOT_AUTH_BANNER"
	botClientBannerEnvName     = "BOT_CLIENT_BANNER"
	botAdvertiserBannerEnvName = "BOT_ADVERTISER_BANNER"
)

// TelegramConfig defines an interface for telegram configuration.
type TelegramConfig interface {
	BotToken() string
	BotAuthBanner() string
	BotClientBanner() string
	BotAdvertiserBanner() string
}

type telegramConfig struct {
	botToken            string
	botAuthBanner       string
	botClientBanner     string
	botAdvertiserBanner string
}

// NewTelegramConfig loads telegram configuration from environment variables.
func NewTelegramConfig() (TelegramConfig, error) {
	botToken := os.Getenv(botTokenEnvName)
	if len(botToken) == 0 {
		return nil, errBotTokenNotSet
	}

	botAuthBanner := os.Getenv(botAuthBannerEnvName)
	if len(botAuthBanner) == 0 {
		return nil, errBotAuthBannerNotSet
	}

	botClientBanner := os.Getenv(botClientBannerEnvName)
	if len(botClientBanner) == 0 {
		return nil, errBotClientBannerNotSet
	}

	botAdvertiserBanner := os.Getenv(botAdvertiserBannerEnvName)
	if len(botAdvertiserBanner) == 0 {
		return nil, errBotAdvertiserBannerNotSet
	}

	return &telegramConfig{
		botToken:            botToken,
		botAuthBanner:       botAuthBanner,
		botClientBanner:     botClientBanner,
		botAdvertiserBanner: botAdvertiserBanner,
	}, nil
}

func (cfg *telegramConfig) BotToken() string {
	return cfg.botToken
}

func (cfg *telegramConfig) BotAuthBanner() string {
	return cfg.botAuthBanner
}

func (cfg *telegramConfig) BotClientBanner() string {
	return cfg.botClientBanner
}

func (cfg *telegramConfig) BotAdvertiserBanner() string {
	return cfg.botAdvertiserBanner
}
