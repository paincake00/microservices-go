package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/inventory/internal/repository"
)

type RepoSuite struct {
	suite.Suite

	partRepo repository.IPartRepository
}

func (s *RepoSuite) SetupSuite() {
	s.partRepo = NewRepository()
}

func (s *RepoSuite) TearDownSuite() {}

func TestRepoIntegration(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
