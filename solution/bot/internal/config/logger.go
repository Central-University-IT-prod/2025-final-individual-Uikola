package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	errIsPrettyNotSet = errors.New("is pretty not set")
	errVersionNotSet  = errors.New("version not set")
	errLogLevelNotSet = errors.New("log level not set")
)

const (
	isPrettyEnvName = "IS_PRETTY"
	versionEnvName  = "VERSION"
	logLevelEnvName = "LOG_LEVEL"
)

// LoggerConfig defines an interface for logger configuration.
type LoggerConfig interface {
	IsPretty() bool
	Version() string
	LogLevel() zerolog.Level
}

type loggerConfig struct {
	isPretty bool
	version  string
	logLevel zerolog.Level
}

// NewLoggerConfig initializes the logger configuration from environment variables.
func NewLoggerConfig() (LoggerConfig, error) {
	isPretty, err := strconv.ParseBool(os.Getenv(isPrettyEnvName))
	if err != nil {
		return nil, errIsPrettyNotSet
	}

	version := os.Getenv(versionEnvName)
	if len(version) == 0 {
		return nil, errVersionNotSet
	}

	logLevel, err := zerolog.ParseLevel(os.Getenv(logLevelEnvName))
	if err != nil {
		return nil, errLogLevelNotSet
	}

	return &loggerConfig{
		isPretty: isPretty,
		version:  version,
		logLevel: logLevel,
	}, nil
}

// IsPretty returns whether pretty logging is enabled.
func (cfg *loggerConfig) IsPretty() bool {
	return cfg.isPretty
}

// Version returns the application version.
func (cfg *loggerConfig) Version() string {
	return cfg.version
}

// LogLevel returns the configured log level.
func (cfg *loggerConfig) LogLevel() zerolog.Level {
	return cfg.logLevel
}
