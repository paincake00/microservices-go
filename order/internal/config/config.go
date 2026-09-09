package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/paincake00/microservices-go/order/internal/config/env"
)

const (
	Path = "CONFIG_PATH"
)

type Config struct {
	Grpc      GrpcConfig
	Http      HttpConfig
	Postgres  PostgresConfig
	Migration MigrationConfig
	Logger    LoggerConfig
}

func Load() (*Config, error) {
	path := fetchConfigPath()

	return LoadFromPath(path)
}

func LoadFromPath(path ...string) (*Config, error) {
	// Выгрузка всех env vars из .env-файлов
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	grpcCfg, err := env.NewGrpcConfig()
	if err != nil {
		return nil, err
	}

	httpCfg, err := env.NewHttpConfig()
	if err != nil {
		return nil, err
	}

	postgresCfg, err := env.NewPostgresCfg()
	if err != nil {
		return nil, err
	}

	migrationCfg, err := env.NewMigrationConfig()
	if err != nil {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Grpc:      grpcCfg,
		Http:      httpCfg,
		Postgres:  postgresCfg,
		Migration: migrationCfg,
		Logger:    loggerCfg,
	}, nil
}

func fetchConfigPath() string {
	return os.Getenv(Path)
}
