package memory

import (
	"testing"
	"time"

	"github.com/paincake00/microservices-go/order/pkg/pgclient"
	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/order/internal/repository"
)

const (
	dbURI = "postgres://postgres:postgres@localhost/postgres?sslmode=disable"

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

	s.orderRepo = NewPostgresOrderStorage(pgClient)
}

func (s *RepoSuite) TearDownSuite() {}

func TestOrderRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
