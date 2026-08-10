package part

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
	"github.com/paincake00/microservices-go/inventory/internal/util/testutil"
)

func (s *RepoSuite) TestCreateAndGetPartSuccess() {
	var (
		uuid = gofakeit.UUID()

		categories, length = testutil.AllCategories()
		category           = categories[rand.Intn(length)]

		expectedPart = testutil.NewPart(uuid, category, testutil.RandomMetadata())
	)

	err := s.partRepo.Create(expectedPart)
	s.NoError(err)

	actualPart, err := s.partRepo.GetPart(uuid)
	s.NoError(err)
	s.Equal(actualPart, expectedPart)
}

func (s *RepoSuite) TestNilPartAndNotFound() {
	uuid := gofakeit.UUID()

	err := s.partRepo.Create(nil)
	s.Error(err)
	s.ErrorIs(err, model.ErrNilPart)

	part, err := s.partRepo.GetPart(uuid)
	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
	s.Nil(part)
}
