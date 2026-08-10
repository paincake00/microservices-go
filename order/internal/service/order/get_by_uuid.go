package order

import (
	"context"

	"github.com/paincake00/microservices-go/order/internal/entity"
)

func (or *Service) GetByUuid(ctx context.Context, orderUUID string) (
	entity.Order,
	error,
) {
	order, err := or.orderStorage.Get(orderUUID)
	if err != nil {
		return entity.Order{}, err
	}
	return order, nil
}
