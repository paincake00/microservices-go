package env

import "github.com/caarlos0/env/v11"

type grpcEnvConfig struct {
	Port string `env:"GRPC_PORT,required"`
}

type GrpcConfig struct {
	raw grpcEnvConfig
}

func NewGrpcConfig() (*GrpcConfig, error) {
	var raw grpcEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}
	return &GrpcConfig{raw: raw}, nil
}

func (gc *GrpcConfig) Port() string {
	return gc.raw.Port
}
