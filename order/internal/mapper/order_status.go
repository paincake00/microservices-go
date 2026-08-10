package mapper

import (
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func OrderStatusToOAS(s enum.OrderStatus) orderv1.OrderStatusEnum {
	switch s {
	case enum.Paid:
		return orderv1.OrderStatusEnumPAID
	case enum.Cancelled:
		return orderv1.OrderStatusEnumCANCELLED
	default:
		return orderv1.OrderStatusEnumPENDINGPAYMENT
	}
}
