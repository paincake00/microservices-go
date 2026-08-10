package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/order/internal/repository"
)

type RepoSuite struct {
	suite.Suite

	orderRepo repository.IOrderRepository
}

func (s *RepoSuite) SetupSuite() {
	s.orderRepo = NewInMemoryStorage()
}

func (s *RepoSuite) TearDownSuite() {}

func TestOrderRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
