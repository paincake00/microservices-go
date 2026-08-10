package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/inventory/internal/repository/mocks"
	"github.com/paincake00/microservices-go/inventory/internal/service"
)

type ServiceSuite struct {
	suite.Suite

	partRepo    *mocks.IPartRepository
	partService service.IPartService
}

func (s *ServiceSuite) SetupSuite() {
	mockRepo := mocks.NewIPartRepository(s.T())

	serv := NewService(mockRepo)

	s.partRepo = mockRepo
	s.partService = serv
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
