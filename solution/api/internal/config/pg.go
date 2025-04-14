package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	errPGUserNotSet = errors.New("pg user not set")
	errPGPassNotSet = errors.New("pg pass not set")
	errPGHostNotSet = errors.New("pg host not set")
	errPGPortNotSet = errors.New("pg port not set")
	errPGDBNotSet   = errors.New("pg db not set")
)

const (
	pgUserEnvName = "PG_USER"
	pgPassEnvName = "PG_PASS"
	pgHostEnvName = "PG_HOST"
	pgPortEnvName = "PG_PORT"
	pgDBEnvName   = "PG_DB"
)

// PGConfig interface defines a method to get the DSN (Data Source Name).
type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	pgUser string
	pgPass string
	pgHost string
	pgPort string
	pgDB   string
}

// NewPGConfig loads pg configuration from environment variables.
func NewPGConfig() (PGConfig, error) {
	pgUser := os.Getenv(pgUserEnvName)
	if len(pgUser) == 0 {
		return nil, errPGUserNotSet
	}
	pgPass := os.Getenv(pgPassEnvName)
	if len(pgPass) == 0 {
		return nil, errPGPassNotSet
	}
	pgHost := os.Getenv(pgHostEnvName)
	if len(pgHost) == 0 {
		return nil, errPGHostNotSet
	}
	pgPort := os.Getenv(pgPortEnvName)
	if len(pgPort) == 0 {
		return nil, errPGPortNotSet
	}
	pgDB := os.Getenv(pgDBEnvName)
	if len(pgDB) == 0 {
		return nil, errPGDBNotSet
	}

	return &pgConfig{
		pgUser: pgUser,
		pgPass: pgPass,
		pgHost: pgHost,
		pgPort: pgPort,
		pgDB:   pgDB,
	}, nil
}

// DSN constructs the pg DSN (Data Source Name).
func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", cfg.pgUser, cfg.pgPass, cfg.pgHost, cfg.pgPort, cfg.pgDB)
}
