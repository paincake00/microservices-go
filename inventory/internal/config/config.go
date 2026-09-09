package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/paincake00/microservices-go/inventory/internal/config/env"
)

const (
	Path = "CONFIG_PATH"
)

type Config struct {
	Grpc   GrpcConfig
	Mongo  MongoConfig
	Logger LoggerConfig
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

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Grpc:   grpcCfg,
		Mongo:  mongoCfg,
		Logger: loggerCfg,
	}, nil
}

func fetchConfigPath() string {
	return os.Getenv(Path)
}
