package memory

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/paincake00/microservices-go/order/internal/repository"
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
	"github.com/paincake00/microservices-go/platform/pkg/migrator"
)

const (
	dbUser = "postgres"
	dbPass = "postgres"
	dbName = "order-db"

	migrationsDir = "../../../../migrations"

	maxOpenCons    = 30
	minIdleCons    = 5
	maxConIdleTime = 5 * time.Second
	maxConLifetime = 30 * time.Second
)

type RepoSuite struct {
	suite.Suite

	container testcontainers.Container
	dbHost    string
	dbPort    string

	orderRepo repository.IOrderRepository
}

func (s *RepoSuite) SetupSuite() {
	ctx := context.Background()

	// Запрос на выкачивание образа Postgres
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.0-alpine3.20",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": dbPass,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForExec(
			[]string{
				"pg_isready",
				"-U", "postgres",
				"-d", "order-db",
				"-h", "localhost",
			},
		).WithStartupTimeout(60 * time.Second),
	}

	// Создание контейнера
	container, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		log.Fatalf("could not start container: %v", err)
	}
	s.container = container

	// Получение хоста контейнера (хост для БД)
	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("could not get host: %v", err)
	}
	s.dbHost = host

	// Получение прокинутого порта для Postgres
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("could not get port for container: %v", err)
	}
	s.dbPort = port.Port()

	// Создание клиента с пулом соединений, настройками и sql-билдером
	pgClient, err := pgclient.New(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, s.dbHost, s.dbPort, dbName),
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
}

func TestOrderRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
