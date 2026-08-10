package service

import (
	"context"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
)

type IOrderService interface {
	Create(ctx context.Context, userUUID string, partUUIDs []string) (
		entity.OrderCompleted,
		error,
	)

	GetByUuid(ctx context.Context, orderUUID string) (
		entity.Order,
		error,
	)

	Pay(
		ctx context.Context,
		orderUUID string,
		payMethod enum.PaymentMethod,
	) (entity.PaymentTransaction, error)

	Cancel(ctx context.Context, orderUUID string) error
}
