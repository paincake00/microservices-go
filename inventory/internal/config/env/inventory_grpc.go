package env

import "github.com/caarlos0/env/v11"

type grpcEnvConfig struct {
	Port string `env:"GRPC_PORT,required"`
}

type GrpcConfig struct {
	raw grpcEnvConfig
}

func NewGrpcConfig() (*GrpcConfig, error) {
	var config grpcEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &GrpcConfig{
		config,
	}, nil
}

func (cfg *GrpcConfig) Port() string {
	return cfg.raw.Port
}
