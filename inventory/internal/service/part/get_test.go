package part

import (
	"context"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		uuid = gofakeit.UUID()

		categories, length = allCategories()

		//nolint:gosec // random is used only for test data generation
		category = categories[rand.Intn(length)]

		part = createFakedPart(uuid, category)

		ctx = context.Background()
	)

	s.partRepo.On("GetPart", ctx, uuid).Return(part, nil).Once()

	res, err := s.partService.GetPart(ctx, uuid)
	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetPartFailure() {
	var (
		uuid = gofakeit.UUID()

		errNotFound = model.ErrNotFound

		ctx = context.Background()
	)

	s.partRepo.On("GetPart", ctx, uuid).Return(nil, errNotFound).Once()

	res, err := s.partService.GetPart(ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, errNotFound)
	s.Nil(res)
}
