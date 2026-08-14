package memory

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/order/internal/migrator"
	"github.com/paincake00/microservices-go/order/internal/repository"
	"github.com/paincake00/microservices-go/order/pkg/pgclient"
)

const (
	dbURI         = "postgres://postgres:postgres@localhost:5432/order-db?sslmode=disable"
	migrationsDir = "../../../../migrations"

	maxOpenCons    = 30
	minIdleCons    = 5
	maxConIdleTime = 5 * time.Second
	maxConLifetime = 30 * time.Second
)

type RepoSuite struct {
	suite.Suite

	orderRepo repository.IOrderRepository
}

func (s *RepoSuite) SetupSuite() {
	// TODO: добавить тестовый Postgres instance через testcontainers

	pgClient, err := pgclient.New(
		dbURI,
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

func (s *RepoSuite) TearDownSuite() {}

func TestOrderRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
