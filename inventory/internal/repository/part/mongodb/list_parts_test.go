package mongodb

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/entity/enum"
	"github.com/paincake00/microservices-go/inventory/internal/util/testutil"
)

func (s *RepoSuite) TestListPartsSuccess() {
	var (
		uuids = []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

		categories = []enum.Category{enum.CategoryFuel, enum.CategoryPorthole, enum.CategoryEngine, enum.CategoryWing}

		// existing parts (UUIDs should be different!)
		parts = []*entity.Part{
			testutil.NewPart(uuids[1], categories[1], testutil.RandomMetadata()),
			testutil.NewPart(uuids[0], categories[2], testutil.RandomMetadata()),
			testutil.NewPart(uuids[2], enum.CategoryUnspecified, testutil.RandomMetadata()), // extra
			testutil.NewPart(gofakeit.UUID(), categories[3], testutil.RandomMetadata()),     // extra
		}

		filter = entity.NewListPartsFilter(uuids, nil, categories, nil, nil)
	)

	// Creating of needed parts
	for _, pp := range parts {
		_, err := s.partRepo.Create(context.Background(), pp)
		s.NoError(err)
	}

	// Testing of filter
	actualParts, err := s.partRepo.ListParts(context.Background(), filter)
	s.NoError(err)
	s.Len(actualParts, 2)
	s.ElementsMatch(
		actualParts, parts[:2],
	) // actualParts может вернуться из map в непредсказуемом порядке. Поэтому проверяем только сами элементы
}
