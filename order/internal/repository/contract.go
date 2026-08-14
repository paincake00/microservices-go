package repository

import (
	"context"

	"github.com/paincake00/microservices-go/order/internal/entity"
)

type IOrderRepository interface {
	Save(ctx context.Context, order entity.Order) (string, error)
	Update(ctx context.Context, order entity.Order) error
	Get(ctx context.Context, orderUUID string) (entity.Order, error)
}
