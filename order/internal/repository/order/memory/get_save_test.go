package memory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
)

func (s *RepoSuite) TestSaveAndGetOrderSuccess() {
	var (
		id = gofakeit.UUID()

		order = entity.Order{
			OrderUuid: id,
			UserUuid:  gofakeit.UUID(),
			PartUuids: []string{
				gofakeit.UUID(),
				gofakeit.UUID(),
			},
			TotalPrice: gofakeit.Float64Range(0, 1000),
			Status:     enum.PendingPayment,
		}
	)

	s.orderRepo.Save(order)

	get, err := s.orderRepo.Get(id)
	s.NoError(err)
	s.Equal(order, get)
}

func (s *RepoSuite) TestGetOrderFailure() {
	id := gofakeit.UUID()

	_, err := s.orderRepo.Get(id)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
