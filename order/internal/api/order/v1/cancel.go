package v1

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func (oh *OrderHandler) CancelOrderByUuid(
	ctx context.Context,
	params orderv1.CancelOrderByUuidParams,
) (orderv1.CancelOrderByUuidRes, error) {
	err := oh.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Cancel order failed: %v", err)

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, entity.ErrOrderConflict) {
			return &orderv1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		}

		return &orderv1.InternalError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderv1.CancelOrderByUuidNoContent{}, nil
}
