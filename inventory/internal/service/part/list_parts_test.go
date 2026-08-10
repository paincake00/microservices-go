package part

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		uuids = []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
			gofakeit.UUID(),
		}

		filter = entity.NewListPartsFilter(uuids, nil, nil, nil, nil)

		categories, length = allCategories()
		//nolint:gosec // random is used only for test data generation
		category = categories[rand.Intn(length)]
		part0    = createFakedPart(uuids[0], category)
		part1    = createFakedPart(uuids[1], category)
		part2    = createFakedPart(uuids[2], category)
	)

	s.partRepo.On("ListParts", filter).Return([]*entity.Part{part0, part1, part2}, nil).Once()

	parts, err := s.partService.ListParts(filter)
	s.NoError(err)
	s.Equal(part0, parts[0])
}

func (s *ServiceSuite) TestListPartsFailure() {
	var (
		filter = entity.NewListPartsFilter(nil, nil, nil, nil, nil)

		errOut = gofakeit.Error()
	)

	s.partRepo.On("ListParts", filter).Return(nil, errOut).Once()

	parts, err := s.partService.ListParts(filter)
	s.Error(err)
	s.ErrorIs(err, errOut)
	s.Nil(parts)
}
