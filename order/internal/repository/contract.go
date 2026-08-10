package repository

import "github.com/paincake00/microservices-go/order/internal/entity"

type IOrderRepository interface {
	Save(order entity.Order)
	Update(order entity.Order)
	Get(orderUUID string) (entity.Order, error)
}
