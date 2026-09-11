//go:build integration

package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/order/internal/repository"
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
	"github.com/paincake00/microservices-go/platform/pkg/logger"
	"github.com/paincake00/microservices-go/platform/pkg/migrator"
	"github.com/paincake00/microservices-go/platform/pkg/testcontainers/postgres"
)

const (
	loggerLevel  = "debug"
	loggerAsJSON = true

	postgresContainer = "test-postgres"
	postgresImage     = "postgres:17.0-alpine3.20"
	postgresUser      = "postgres"
	postgresPass      = "postgres"
	postgresDatabase  = "order-db"

	migrationsDir = "../../../../migrations"

	maxOpenCons    = 30
	minIdleCons    = 5
	maxConIdleTime = 5 * time.Second
	maxConLifetime = 30 * time.Second
)

type RepoSuite struct {
	suite.Suite

	container *postgres.Container

	orderRepo repository.IOrderRepository
}

func (s *RepoSuite) SetupSuite() {
	ctx := context.Background()

	// Инициализаия логгера
	logger.Init(loggerLevel, loggerAsJSON)

	// Создание Postgres
	container, err := postgres.NewContainer(
		ctx,
		postgres.WithNetworkName(""), // сеть не нужна
		postgres.WithContainerName(postgresContainer),
		postgres.WithImageName(postgresImage),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithUsername(postgresUser),
		postgres.WithPassword(postgresPass),
		postgres.WithLogger(logger.Logger()),
	)
	s.Require().NoError(err)

	s.container = container

	// Создание клиента с пулом соединений, настройками и sql-билдером
	pgClient, err := pgclient.New(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPass, s.container.Host(),
			s.container.Port(), postgresDatabase,
		),
		pgclient.MaxOpenCons(maxOpenCons),
		pgclient.MinIdleCons(minIdleCons),
		pgclient.MaxConIdleTime(maxConIdleTime),
		pgclient.MaxConLifetime(maxConLifetime),
	)
	s.Require().NoError(err)

	migration := migrator.NewMigrator(stdlib.OpenDBFromPool(pgClient.Pool), migrationsDir)

	// Применяем миграцию (если уже была, то ничего не будет, Goose умный)
	err = migration.Up()
	s.Require().NoError(err)

	s.orderRepo = NewPostgresOrderStorage(pgClient)
}

func (s *RepoSuite) TearDownSuite() {
	err := s.container.Terminate(context.Background())
	s.Require().NoError(err)
	logger.Info(context.Background(), "🛑 Контейнер Postgres остановлен")
}

func TestOrderRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
