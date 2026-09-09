package mongo

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	mongoPort           = "27017"
	mongoStartupTimeout = 1 * time.Minute

	mongoEnvUsernameKey = "MONGO_INITDB_ROOT_USERNAME"
	mongoEnvPasswordKey = "MONGO_INITDB_ROOT_PASSWORD"
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
				cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Mongo container started")
	success = true

	return &Container{container: container, cfg: cfg}, nil
}

func (c *Container) Host() string {
	return c.cfg.Host
}

func (c *Container) Port() string {
	return c.cfg.Port
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Mongo container terminated")

	return nil
}

func initContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			mongoEnvUsernameKey: cfg.Username,
			mongoEnvPasswordKey: cfg.Password,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(mongoPort+"/tcp"),
			wait.ForExec(
				[]string{
					"mongosh",
					"--quiet",
					"-u", cfg.Username,
					"-p", cfg.Password,
					"--authenticationDatabase", cfg.AuthDB,
					"--eval", "db.runCommand({ ping: 1 }).ok",
				},
			),
		).WithDeadline(mongoStartupTimeout),
		//HostConfigModifier: defaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, errors.Errorf("failed to start mongo container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, mongoPort+"/tcp")
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, port.Port(), nil
}
