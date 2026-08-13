package memory

import (
	"github.com/paincake00/microservices-go/order/internal/entity"
)

func (o *PostgresOrderStorage) Save(order entity.Order) {
	o.pgClient.Pool
}

func (o *PostgresOrderStorage) Update(order entity.Order) {
	o.Save(order)
}
