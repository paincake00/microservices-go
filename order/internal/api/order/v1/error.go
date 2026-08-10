package v1

import (
	"context"
	"net/http"

	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

// NewError создает новую ошибку в формате GenericError
func (oh *OrderHandler) NewError(_ context.Context, err error) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderv1.GenericError{
			Code:    orderv1.NewOptInt(http.StatusInternalServerError),
			Message: orderv1.NewOptString(err.Error()),
		},
	}
}
