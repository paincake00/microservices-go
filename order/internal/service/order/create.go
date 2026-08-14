package order

import (
	"context"
	"log"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
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

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	order := entity.Order{
		UserUuid:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     enum.PendingPayment,
	}

	orderUuid, err := or.orderStorage.Save(ctx, order)
	if err != nil {
		return entity.OrderCompleted{}, err
	}

	log.Printf("Create order: %v", orderUuid)

	return entity.OrderCompleted{
		OrderUuid:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}
