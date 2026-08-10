package part

import (
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
	)

	s.partRepo.On("GetPart", uuid).Return(part, nil).Once()

	res, err := s.partService.GetPart(uuid)
	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetPartFailure() {
	var (
		uuid = gofakeit.UUID()

		errNotFound = model.ErrNotFound
	)

	s.partRepo.On("GetPart", uuid).Return(nil, errNotFound).Once()

	res, err := s.partService.GetPart(uuid)
	s.Error(err)
	s.ErrorIs(err, errNotFound)
	s.Nil(res)
}
