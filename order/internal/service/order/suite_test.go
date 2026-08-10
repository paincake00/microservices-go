package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	mocksClient "github.com/paincake00/microservices-go/order/internal/client/grpc/mocks"
	mocksRepo "github.com/paincake00/microservices-go/order/internal/repository/mocks"
	"github.com/paincake00/microservices-go/order/internal/service"
)

type ServiceSuite struct {
	suite.Suite

	// nolint:containedctx // for tests
	ctx context.Context

	orderRepo        *mocksRepo.IOrderRepository
	inventoryService *mocksClient.InventoryService
	paymentService   *mocksClient.PaymentService

	orderService service.IOrderService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.orderRepo = mocksRepo.NewIOrderRepository(s.T())
	s.inventoryService = mocksClient.NewInventoryService(s.T())
	s.paymentService = mocksClient.NewPaymentService(s.T())

	s.orderService = NewOrderService(s.orderRepo, s.inventoryService, s.paymentService)
}

func (s *ServiceSuite) TearDownSuite() {}

func TestOrderRepoIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
