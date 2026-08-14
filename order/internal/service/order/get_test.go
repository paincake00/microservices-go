package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
)

func (s *ServiceSuite) TestGetByUuidFailure() {
	var (
		orderUuid = gofakeit.UUID()

		errNotFound = model.ErrOrderNotFound
	)

	// настройка моков
	s.orderRepo.On("Get", context.Background(), orderUuid).Return(entity.Order{}, errNotFound).Once()

	order, err := s.orderService.GetByUuid(s.ctx, orderUuid)
	s.Error(err)
	s.ErrorIs(err, errNotFound)
	s.Empty(order)
}
