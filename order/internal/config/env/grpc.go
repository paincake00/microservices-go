package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	InventoryHost      string        `env:"GRPC_INVENTORY_HOST,required"`
	InventoryPort      string        `env:"GRPC_INVENTORY_PORT,required"`
	PaymentHost        string        `env:"GRPC_PAYMENT_HOST,required"`
	PaymentPort        string        `env:"GRPC_PAYMENT_PORT,required"`
	HealthCheckTimeout time.Duration `env:"GRPC_HEALTH_CHECK_TIMEOUT,required"`
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

func (cfg *GrpcConfig) GetInventoryHost() string {
	return cfg.raw.InventoryHost
}

func (cfg *GrpcConfig) GetInventoryPort() string {
	return cfg.raw.InventoryPort
}

func (cfg *GrpcConfig) GetPaymentHost() string {
	return cfg.raw.PaymentHost
}

func (cfg *GrpcConfig) GetPaymentPort() string {
	return cfg.raw.PaymentPort
}

func (cfg *GrpcConfig) GetHealthCheckTimeout() time.Duration {
	return cfg.raw.HealthCheckTimeout
}
