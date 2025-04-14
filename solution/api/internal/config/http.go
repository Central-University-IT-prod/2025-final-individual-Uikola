package config

import (
	"errors"
	"net"
	"os"
)

var (
	errHTTPHostNotSet = errors.New("http host not set")
	errHTTPPortNotSet = errors.New("http port not set")
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

// HTTPConfig defines an interface for HTTP server configuration.
type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig initializes a new HTTP configuration from environment variables.
func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errHTTPHostNotSet
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errHTTPPortNotSet
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

// Address constructs and returns the full server address (host:port).
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
