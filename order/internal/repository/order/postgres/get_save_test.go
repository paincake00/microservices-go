//go:build integration

package memory

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
)

func (s *RepoSuite) TestSaveAndGetOrderSuccess() {
	order := entity.Order{
		UserUuid: gofakeit.UUID(),
		PartUuids: []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
		},
		TotalPrice: gofakeit.Float64Range(0, 1000),
		Status:     enum.PendingPayment,
	}

	orderUuid, err := s.orderRepo.Save(context.Background(), order)
	s.NoError(err)

	get, err := s.orderRepo.Get(context.Background(), orderUuid)
	s.NoError(err)
	order.OrderUuid = orderUuid
	s.Equal(order, get)
}

func (s *RepoSuite) TestGetOrderFailure() {
	id := gofakeit.UUID()

	_, err := s.orderRepo.Get(context.Background(), id)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
