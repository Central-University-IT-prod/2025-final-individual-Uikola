package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	errClickHouseHostNotSet = errors.New("clickhouse host not set")
	errClickHousePortNotSet = errors.New("clickhouse port not set")
	errClickHouseUserNotSet = errors.New("clickhouse user not set")
	errClickHousePassNotSet = errors.New("clickhouse password not set")
	errClickHouseDBNotSet   = errors.New("clickhouse database not set")
)

const (
	clickHouseHostEnvName = "CLICK_HOUSE_HOST"
	clickHousePortEnvName = "CLICK_HOUSE_PORT"
	clickHouseUserEnvName = "CLICK_HOUSE_USER"
	clickHousePassEnvName = "CLICK_HOUSE_PASS"
	clickHouseDBEnvName   = "CLICK_HOUSE_DB"
)

// ClickHouseConfig defines an interface for Clickhouse server.
type ClickHouseConfig interface {
	DSN() string
}

type clickHouseConfig struct {
	clickHouseHost string
	clickHousePort string
	clickHouseUser string
	clickHousePass string
	clickHouseDB   string
}

// NewClickHouseConfig initializes a new ClickHouse configuration from environment variables.
func NewClickHouseConfig() (ClickHouseConfig, error) {
	clickHouseHost := os.Getenv(clickHouseHostEnvName)
	if len(clickHouseHost) == 0 {
		return nil, errClickHouseHostNotSet
	}

	clickHousePort := os.Getenv(clickHousePortEnvName)
	if len(clickHousePort) == 0 {
		return nil, errClickHousePortNotSet
	}

	clickHouseUser := os.Getenv(clickHouseUserEnvName)
	if len(clickHouseUser) == 0 {
		return nil, errClickHouseUserNotSet
	}

	clickHousePass := os.Getenv(clickHousePassEnvName)
	if len(clickHousePass) == 0 {
		return nil, errClickHousePassNotSet
	}

	clickHouseDB := os.Getenv(clickHouseDBEnvName)
	if len(clickHouseDB) == 0 {
		return nil, errClickHouseDBNotSet
	}

	return &clickHouseConfig{
		clickHouseHost: clickHouseHost,
		clickHousePort: clickHousePort,
		clickHouseUser: clickHouseUser,
		clickHousePass: clickHousePass,
		clickHouseDB:   clickHouseDB,
	}, nil
}

// DSN constructs the clickhouse DSN (Data Source Name).
func (cfg *clickHouseConfig) DSN() string {
	return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", cfg.clickHouseUser, cfg.clickHousePass, cfg.clickHouseHost, cfg.clickHousePort, cfg.clickHouseDB)
	//return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s", cfg.clickHouseHost, cfg.clickHousePort, cfg.clickHouseDB, cfg.clickHouseUser, cfg.clickHousePass)
}
