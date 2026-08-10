package mapper

import (
	"github.com/paincake00/microservices-go/order/internal/entity"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func OrderToOAS(o entity.Order) *orderv1.OrderByUUIDResponse {
	order := orderv1.OrderByUUIDResponse{
		OrderUUID:  o.OrderUuid,
		UserUUID:   o.UserUuid,
		PartUuids:  o.PartUuids,
		TotalPrice: o.TotalPrice,
		Status:     OrderStatusToOAS(o.Status),
	}

	if o.TransactionUuid != nil {
		order.TransactionUUID.SetTo(*o.TransactionUuid)
	}
	if o.PaymentMethod != nil {
		order.PaymentMethod.SetTo(PaymentMethodToOAS(*o.PaymentMethod))
	}

	return &order
}
