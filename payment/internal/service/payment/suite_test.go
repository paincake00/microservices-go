package payment

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/paincake00/microservices-go/payment/internal/service"
	"github.com/paincake00/microservices-go/payment/internal/util/uuidutil"
)

type UUIDGeneratorStub struct{}

func (s *UUIDGeneratorStub) NewV7() (uuid.UUID, error) {
	return uuid.Nil, errors.New("generation failed")
}

type ServiceSuite struct {
	suite.Suite

	payService service.IPaymentService

	corruptedPayService service.IPaymentService
}

func (s *ServiceSuite) SetupSuite() {
	s.payService = NewService(
		uuidutil.NewUuidGeneratorV7(),
	)

	s.corruptedPayService = NewService(
		new(UUIDGeneratorStub),
	)
}

func (s *ServiceSuite) TearDownSuite() {}

func TestPayServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
