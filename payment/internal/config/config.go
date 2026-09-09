package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/paincake00/microservices-go/payment/internal/config/env"
)

const (
	path = "CONFIG_PATH"
)

type Config struct {
	Grpc   GrpcConfig
	Logger LoggerConfig
}

func Load() (*Config, error) {
	cfgPath := fetchConfigPath()

	return LoadFromPath(cfgPath)
}

func LoadFromPath(path ...string) (*Config, error) {
	err := godotenv.Load(path...)
	if err != nil {
		return nil, err
	}

	grpcCfg, err := env.NewGrpcConfig()
	if err != nil {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Grpc:   grpcCfg,
		Logger: loggerCfg,
	}, nil
}

func fetchConfigPath() string {
	return os.Getenv(path)
}
