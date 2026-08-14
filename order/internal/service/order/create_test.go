package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
)

func (s *ServiceSuite) TestCreateAndGetOrderSuccess() {
	var (
		userUuid = gofakeit.UUID()

		partUuids = []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
			gofakeit.UUID(),
		}

		totalPrice0   = gofakeit.Float64()
		totalPrice1   = gofakeit.Float64()
		totalPrice2   = gofakeit.Float64()
		expectedParts = []*entity.Part{
			{
				Price: totalPrice0,
			},
			{
				Price: totalPrice1,
			},
			{
				Price: totalPrice2,
			},
		}
	)

	// настройка моков
	s.orderRepo.On("Save", context.Background(), mock.Anything).Once()
	s.inventoryService.On("ListParts", context.Background(), partUuids).Return(expectedParts, nil).Once()

	orderResp, err := s.orderService.Create(s.ctx, userUuid, partUuids)
	s.NoError(err)
	s.NotEmpty(orderResp)

	order := entity.Order{
		OrderUuid:  orderResp.OrderUuid,
		UserUuid:   userUuid,
		PartUuids:  partUuids,
		TotalPrice: totalPrice0 + totalPrice1 + totalPrice2,
		Status:     enum.PendingPayment,
	}

	s.orderRepo.On("Get", context.Background(), orderResp.OrderUuid).Return(order, nil).Once()

	actualOrder, err := s.orderService.GetByUuid(s.ctx, orderResp.OrderUuid)
	s.NoError(err)
	s.Equal(order, actualOrder)
}

func (s *ServiceSuite) TestCreateFailure() {
	var (
		userUuid = gofakeit.UUID()

		partUuids = []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
			gofakeit.UUID(),
		}

		expectedParts = []*entity.Part{
			{
				Price: gofakeit.Float64Range(0, 1000),
			},
			{
				Price: gofakeit.Float64Range(0, 1000),
			},
			{
				Price: gofakeit.Float64Range(0, 1000),
			},
		}
	)

	// настройка моков
	s.orderRepo.On("Save", context.Background(), mock.Anything).Maybe()
	s.inventoryService.On("ListParts", context.Background(), partUuids).Return(expectedParts[:2], nil).Once()

	orderResp, err := s.orderService.Create(s.ctx, userUuid, partUuids)
	s.Error(err)
	s.ErrorIs(err, entity.ErrPartsNotFound)
	s.Empty(orderResp)
}
