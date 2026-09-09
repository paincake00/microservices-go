//go:build integration

package mongodb

import (
	"context"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
	"github.com/paincake00/microservices-go/inventory/internal/util/testutil"
)

func (s *RepoSuite) TestCreateAndGetPartSuccess() {
	var (
		uuid = gofakeit.UUID()

		categories, length = testutil.AllCategories()
		category           = categories[rand.Intn(length)]

		expectedPart = testutil.NewPart(uuid, category, testutil.RandomMetadata())

		ctx = context.Background()
	)

	id, err := s.partRepo.Create(ctx, expectedPart)
	s.NoError(err)

	actualPart, err := s.partRepo.GetPart(ctx, id)
	s.NoError(err)
	s.Equal(expectedPart, actualPart)
}

func (s *RepoSuite) TestNilPartAndNotFound() {
	ctx := context.Background()

	id, err := s.partRepo.Create(ctx, nil)
	s.Error(err)
	s.ErrorIs(err, model.ErrNilPart)

	part, err := s.partRepo.GetPart(ctx, id)
	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
	s.Nil(part)
}

func (s *RepoSuite) TestValueBSON() {
	original := entity.NewInt64Value(123)

	data, err := bson.Marshal(original)
	s.NoError(err)

	var decoded entity.Value
	err = bson.Unmarshal(data, &decoded)
	s.NoError(err)

	s.Equal(original.GetKind(), decoded.GetKind())
	s.Equal(original.GetInt64Value(), decoded.GetInt64Value())
}
