package v1

import "github.com/paincake00/microservices-go/order/internal/service"

type OrderHandler struct {
	orderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}
