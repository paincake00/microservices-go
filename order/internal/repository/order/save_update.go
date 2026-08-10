package order

import (
	"github.com/paincake00/microservices-go/order/internal/entity"
)

func (o *InMemoryStorage) Save(order entity.Order) {
	o.mx.Lock()
	defer o.mx.Unlock()

	o.orders[order.OrderUuid] = order
}

func (o *InMemoryStorage) Update(order entity.Order) {
	o.Save(order)
}
