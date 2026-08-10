package grpc

import (
	"context"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
)

type InventoryService interface {
	ListParts(ctx context.Context, partUuids []string) ([]*entity.Part, error)
}

type PaymentService interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod enum.PaymentMethod) (
		entity.PaymentTransaction,
		error,
	)
}
