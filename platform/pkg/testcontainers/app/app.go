package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/go-faster/errors"
	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/paincake00/microservices-go/platform/pkg/logger"
)

const (
	defaultAppName        = "app"
	defaultAppPort        = "50051"
	defaultStartupTimeout = 1 * time.Minute
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	Name          string
	DockerfileDir string
	Dockerfile    string
	Port          string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Binds         []Mount
	Logger        Logger
}

type Mount struct {
	Source   string
	Target   string
	ReadOnly bool
}

type Container struct {
	container    testcontainers.Container
	externalHost string
	externalPort string
	cfg          *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := &Config{
		Name:          defaultAppName,
		Port:          defaultAppPort,
		Dockerfile:    "Dockerfile",
		DockerfileDir: ".",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultAppPort + "/tcp").WithStartupTimeout(defaultStartupTimeout),
		Env:           make(map[string]string),
		Binds:         []Mount{},
		Logger:        &logger.NoopLogger{},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Name: cfg.Name,
		FromDockerfile: testcontainers.FromDockerfile{
			Context:        cfg.DockerfileDir,
			Dockerfile:     cfg.Dockerfile,
			BuildLogWriter: cfg.LogOutput,
		},
		Networks:           cfg.Networks,
		Env:                cfg.Env,
		WaitingFor:         cfg.StartupWait,
		ExposedPorts:       []string{cfg.Port + "/tcp"},
		HostConfigModifier: addAllBinds(cfg.Binds),
	}

	genericContainer, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, errors.Errorf("failed to start app genericContainer: %v", err)
	}

	mappedPort, err := genericContainer.MappedPort(ctx, cfg.Port+"/tcp")
	if err != nil {
		return nil, errors.Errorf("failed to get mapped externalPort: %v", err)
	}

	host, err := genericContainer.Host(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get genericContainer externalHost: %v", err)
	}

	cfg.Logger.Info(ctx, "App container started", zap.String("uri:", net.JoinHostPort(host, mappedPort.Port())))

	return &Container{
		container:    genericContainer,
		externalHost: host,
		externalPort: mappedPort.Port(),
		cfg:          cfg,
	}, nil
}

func (a *Container) Address() string {
	return net.JoinHostPort(a.externalHost, a.externalPort)
}

func (a *Container) Terminate(ctx context.Context) error {
	return a.container.Terminate(ctx)
}

func addAllBinds(binds []Mount) func(hostConfig *container.HostConfig) {
	return func(hostConfig *container.HostConfig) {
		for _, bind := range binds {
			mode := "rw"
			if bind.ReadOnly {
				mode = "ro"
			}

			hostConfig.Binds = append(
				hostConfig.Binds,
				fmt.Sprintf("%s:%s:%s", bind.Source, bind.Target, mode),
			)
		}
	}
}
