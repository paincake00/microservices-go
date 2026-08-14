package order

import (
	"context"
	"log"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
)

func (or *Service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := or.orderStorage.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.Status == enum.Paid {
		return entity.ErrOrderConflict
	}

	if order.Status == enum.PendingPayment {
		order.Status = enum.Cancelled
		err = or.orderStorage.Update(ctx, order)
		if err != nil {
			return err
		}
		log.Printf("order cancelled with UUID: %s", orderUUID)
	}
	return nil
}
