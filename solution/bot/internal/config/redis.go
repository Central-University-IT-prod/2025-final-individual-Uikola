package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	errRedisHostNotSet = errors.New("redis host not set")
	errRedisPortNotSet = errors.New("redis port not set")
	errRedisPassNotSet = errors.New("redis pass not set")
)

const (
	redisHostEnvName = "REDIS_HOST"
	redisPortEnvName = "REDIS_PORT"
	redisPassEnvName = "REDIS_PASS"
)

// RedisConfig defines an interface for redis configuration.
type RedisConfig interface {
	Addr() string
	Pass() string
}

type redisConfig struct {
	host string
	port string
	pass string
}

// NewRedisConfig initializes the redis configuration from environment variables.
func NewRedisConfig() (RedisConfig, error) {
	redisHost := os.Getenv(redisHostEnvName)
	if len(redisHost) == 0 {
		return nil, errRedisHostNotSet
	}

	redisPort := os.Getenv(redisPortEnvName)
	if len(redisPort) == 0 {
		return nil, errRedisPortNotSet
	}

	redisPass := os.Getenv(redisPassEnvName)
	if len(redisPass) == 0 {
		return nil, errRedisPassNotSet
	}

	return &redisConfig{
		host: redisHost,
		port: redisPort,
		pass: redisPass,
	}, nil
}

// Addr returns redis host.
func (cfg *redisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", cfg.host, cfg.port)
}

// Pass returns redis password.
func (cfg *redisConfig) Pass() string {
	return cfg.pass
}
