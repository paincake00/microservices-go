package memory

import (
	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
)

func (o *InMemoryOrderStorage) Get(orderUUID string) (entity.Order, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	if order, ok := o.orders[orderUUID]; ok {
		return order, nil
	}
	return entity.Order{}, model.ErrOrderNotFound
}
