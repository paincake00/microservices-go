package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	var (
		orderID = gofakeit.UUID()

		order = entity.Order{
			Status: enum.PendingPayment,
		}

		ctx = context.Background()
	)

	s.orderRepo.On("Get", ctx, orderID).Return(order, nil).Once()
	order.Status = enum.Cancelled
	s.orderRepo.On("Update", ctx, order).Return(nil).Once()

	err := s.orderService.Cancel(ctx, orderID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderFailure() {
	var (
		orderID = gofakeit.UUID()

		order = entity.Order{
			Status: enum.Paid,
		}

		ctx = context.Background()

		errConflict = entity.ErrOrderConflict
	)

	s.orderRepo.On("Get", ctx, orderID).Return(order, nil).Once()

	err := s.orderService.Cancel(ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, errConflict)
}
