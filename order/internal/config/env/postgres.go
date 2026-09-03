package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	Db       string `env:"POSTGRES_DB,required"`

	MaxOpenCons    int           `env:"POSTGRES_MAX_OPEN_CONS,required"`
	MinIdleCons    int           `env:"POSTGRES_MIN_IDLE_CONS,required"`
	MaxConIdleTime time.Duration `env:"POSTGRES_MAX_CON_IDLE_TIME,required"`
	MaxConLifetime time.Duration `env:"POSTGRES_MAX_CON_LIFETIME,required"`
}

type PostgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresCfg() (*PostgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &PostgresConfig{raw: raw}, nil
}

func (cfg *PostgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.raw.User, cfg.raw.Password, cfg.raw.Host, cfg.raw.Port, cfg.raw.Db,
	)
}

func (cfg *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.raw.Host, cfg.raw.Port, cfg.raw.User, cfg.raw.Password, cfg.raw.Db,
	)
}

func (cfg *PostgresConfig) MaxOpenCons() int {
	return cfg.raw.MaxOpenCons
}

func (cfg *PostgresConfig) MinIdleCons() int {
	return cfg.raw.MinIdleCons
}

func (cfg *PostgresConfig) MaxConIdleTime() time.Duration {
	return cfg.raw.MaxConIdleTime
}

func (cfg *PostgresConfig) MaxConLifetime() time.Duration {
	return cfg.raw.MaxConLifetime
}
