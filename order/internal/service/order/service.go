package order

import (
	"github.com/paincake00/microservices-go/order/internal/client/grpc"
	"github.com/paincake00/microservices-go/order/internal/repository"
)

type Service struct {
	orderStorage     repository.IOrderRepository
	inventoryService grpc.InventoryService
	paymentService   grpc.PaymentService
}

func NewOrderService(
	orderStorage repository.IOrderRepository,
	inventoryService grpc.InventoryService,
	paymentService grpc.PaymentService,
) *Service {
	return &Service{
		orderStorage:     orderStorage,
		inventoryService: inventoryService,
		paymentService:   paymentService,
	}
}
