package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/paincake00/microservices-go/order/internal/mapper"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func (oh *OrderHandler) GetOrderByUuid(
	ctx context.Context,
	params orderv1.GetOrderByUuidParams,
) (orderv1.GetOrderByUuidRes, error) {
	order, err := oh.orderService.GetByUuid(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
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

	return mapper.OrderToOAS(order), nil
}
