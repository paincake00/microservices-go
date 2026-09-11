package postgres

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	postgresPort           = "5432"
	postgresStartupTimeout = 1 * time.Minute

	postgresEnvUsernameKey = "POSTGRES_USER"
	//nolint:gosec
	postgresEnvPasswordKey = "POSTGRES_PASSWORD"
	postgresEnvDBKey       = "POSTGRES_DB"
)

type Container struct {
	container testcontainers.Container
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := initContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Postgres container started")
	success = true

	return &Container{container: container, cfg: cfg}, nil
}

func (c *Container) Host() string {
	return c.cfg.Host
}

func (c *Container) Port() string {
	return c.cfg.Port
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Postgres container terminated")

	return nil
}

func initContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			postgresEnvUsernameKey: cfg.Username,
			postgresEnvPasswordKey: cfg.Password,
			postgresEnvDBKey:       cfg.Database,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(postgresPort+"/tcp"),
			wait.ForExec(
				[]string{
					"pg_isready",
					"-U", cfg.Username,
					"-d", cfg.Database,
					"-h", "localhost",
				},
			),
		).WithDeadline(postgresStartupTimeout),
	}

	container, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, errors.Errorf("failed to start postgres container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, postgresPort+"/tcp")
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, port.Port(), nil
}
