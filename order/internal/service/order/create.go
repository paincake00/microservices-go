package order

import (
	"context"
	"log"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
	"github.com/paincake00/microservices-go/order/internal/util/utiluuid"
)

// Create returns OrderUUID and TotalPrice of all parts.
// Otherwise, it returns error
func (or *Service) Create(ctx context.Context, userUUID string, partUUIDs []string) (
	entity.OrderCompleted,
	error,
) {
	parts, err := or.inventoryService.ListParts(ctx, partUUIDs)
	if err != nil {
		return entity.OrderCompleted{}, err
	}
	if len(parts) < len(partUUIDs) {
		return entity.OrderCompleted{}, entity.ErrPartsNotFound
	}

	orderUuid, err := utiluuid.GetNewUUID()
	if err != nil {
		return entity.OrderCompleted{}, err
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	order := entity.Order{
		OrderUuid:  orderUuid,
		UserUuid:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     enum.PendingPayment,
	}

	log.Printf("Create order: %+v", order)

	or.orderStorage.Save(order)

	return entity.OrderCompleted{
		OrderUuid:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}
