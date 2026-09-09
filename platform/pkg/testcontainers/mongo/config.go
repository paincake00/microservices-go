package mongo

import (
	"context"

	"github.com/moby/moby/api/types/container"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName   string
	ImageName     string
	ContainerName string
	Username      string
	Password      string
	Database      string
	AuthDB        string
	Logger        Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "",
		ContainerName: "mongo-container",
		ImageName:     "mongo:8.0",
		Database:      "test",
		Username:      "root",
		Password:      "root",
		AuthDB:        "admin",
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
