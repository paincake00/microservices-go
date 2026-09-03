package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_DATABASE,required"`
	User     string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
}

type MongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*MongoConfig, error) {
	var config mongoEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &MongoConfig{config}, nil
}

func (cfg *MongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.User, cfg.raw.Password, cfg.raw.Host, cfg.raw.Port, cfg.raw.Database, cfg.raw.AuthDB,
	)
}

func (cfg *MongoConfig) DatabaseName() string {
	return cfg.raw.Database
}
