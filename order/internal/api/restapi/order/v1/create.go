package v1

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/paincake00/microservices-go/order/internal/entity"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func (oh *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (
	orderv1.CreateOrderRes,
	error,
) {
	resp, err := oh.orderService.Create(ctx, req.GetUserUUID(), req.GetPartUuids())
	if err != nil {
		log.Printf("Create order failed: %v", err)

		if errors.Is(err, entity.ErrPartsNotFound) {
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		return &orderv1.InternalError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  resp.OrderUuid,
		TotalPrice: resp.TotalPrice,
	}, nil
}
