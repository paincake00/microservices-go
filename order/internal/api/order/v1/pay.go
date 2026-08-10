package v1

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/paincake00/microservices-go/order/internal/mapper"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func (oh *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderv1.PayOrderRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	responseWithTransaction, err := oh.orderService.Pay(
		ctx, params.OrderUUID, mapper.PaymentMethodFromOAS(req.GetPaymentMethod()),
	)
	if err != nil {
		log.Printf("Pay order failed: %v", err)

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

	return &orderv1.PayOrderResponse{
		TransactionUUID: responseWithTransaction.TransactionUuid,
	}, nil
}
