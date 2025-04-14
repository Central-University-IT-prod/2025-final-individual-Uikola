package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	errMinioHostNotSet       = errors.New("minio host not set")
	errMinioPortNotSet       = errors.New("minio port not set")
	errMinioPublicHostNotSet = errors.New("minio public host not set")
	errMinioPublicPortNotSet = errors.New("minio public port not set")
	errMinioUserNotSet       = errors.New("minio user not set")
	errMinioPassNotSet       = errors.New("minio password not set")
	errMinioUseSSLNotSet     = errors.New("minio use ssl not set")
	errMinioBucketNameNotSet = errors.New("minio bucket name not set")
)

const (
	minioHostEnvName       = "MINIO_HOST"
	minioPortEnvName       = "MINIO_PORT"
	minioPublicHostEnvName = "MINIO_PUBLIC_HOST"
	minioPublicPortEnvName = "MINIO_PUBLIC_PORT"
	minioUserEnvName       = "MINIO_USER"
	minioPassEnvName       = "MINIO_PASS"
	minioUseSSLEnvName     = "MINIO_USE_SSL"
	minioBucketNameEnvName = "MINIO_BUCKET_NAME"
)

// MinioConfig defines an interface for minio s3 config.
type MinioConfig interface {
	Endpoint() string
	PublicEndpoint() string
	User() string
	Pass() string
	UseSSL() bool
	BucketName() string
}

type minioConfig struct {
	minioHost       string
	minioPort       string
	minioPublicHost string
	minioPublicPort string
	minioUser       string
	minioPass       string
	minioUseSSL     bool
	minioBucketName string
}

// NewMinioConfig initializes a new Minio configuration from environment variables.
func NewMinioConfig() (MinioConfig, error) {
	minioHost := os.Getenv(minioHostEnvName)
	if len(minioHost) == 0 {
		return nil, errMinioHostNotSet
	}

	minioPort := os.Getenv(minioPortEnvName)
	if len(minioPort) == 0 {
		return nil, errMinioPortNotSet
	}

	minioPublicHost := os.Getenv(minioPublicHostEnvName)
	if len(minioPublicHost) == 0 {
		return nil, errMinioPublicHostNotSet
	}

	minioPublicPort := os.Getenv(minioPublicPortEnvName)
	if len(minioPublicPort) == 0 {
		return nil, errMinioPublicPortNotSet
	}

	minioUser := os.Getenv(minioUserEnvName)
	if len(minioUser) == 0 {
		return nil, errMinioUserNotSet
	}

	minioPass := os.Getenv(minioPassEnvName)
	if len(minioPass) == 0 {
		return nil, errMinioPassNotSet
	}

	minioUseSSL, err := strconv.ParseBool(os.Getenv(minioUseSSLEnvName))
	if err != nil {
		return nil, errMinioUseSSLNotSet
	}

	minioBucketName := os.Getenv(minioBucketNameEnvName)
	if len(minioBucketName) == 0 {
		return nil, errMinioBucketNameNotSet
	}

	return &minioConfig{
		minioHost:       minioHost,
		minioPort:       minioPort,
		minioPublicHost: minioPublicHost,
		minioPublicPort: minioPublicPort,
		minioUser:       minioUser,
		minioPass:       minioPass,
		minioUseSSL:     minioUseSSL,
		minioBucketName: minioBucketName,
	}, nil
}

// Endpoint constructs and returns the full endpoint URL for the MinIO server.
func (cfg *minioConfig) Endpoint() string {
	return fmt.Sprintf("%s:%s", cfg.minioHost, cfg.minioPort)
}

// PublicEndpoint constructs and returns the full public endpoint URL for the MinIO server.
func (cfg *minioConfig) PublicEndpoint() string {
	return fmt.Sprintf("%s:%s", cfg.minioPublicHost, cfg.minioPublicPort)
}

// User returns the username for MinIO authentication.
func (cfg *minioConfig) User() string {
	return cfg.minioUser
}

// Pass returns the password for MinIO authentication.
func (cfg *minioConfig) Pass() string {
	return cfg.minioPass
}

// UseSSL indicates whether SSL is enabled for the MinIO connection.
func (cfg *minioConfig) UseSSL() bool {
	return cfg.minioUseSSL
}

// BucketName returns the name of the S3 bucket being used.
func (cfg *minioConfig) BucketName() string {
	return cfg.minioBucketName
}
