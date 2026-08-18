package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/mapper"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

func (p *PartStorage) ListParts(ctx context.Context, filter entity.ListPartsFilter) ([]*entity.Part, error) {
	parts, err := p.collection.Find(ctx, buildBSONFilter(filter))
	if err != nil {
		return nil, err
	}
	// nolint:errcheck
	defer parts.Close(ctx)

	res := make([]*entity.Part, 0)
	for parts.Next(ctx) {
		var partModel model.Part

		err = parts.Decode(&partModel)
		if err != nil {
			return nil, err
		}

		part, errMap := mapper.FromModel(&partModel)
		if errMap != nil {
			return nil, errMap
		}

		res = append(res, part)
	}

	if err = parts.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func buildBSONFilter(filter entity.ListPartsFilter) bson.D {
	bsonFilter := make(bson.A, 0, 5)

	if len(filter.Uuids) > 0 {
		bsonFilter = append(
			bsonFilter, bson.D{
				{
					Key: "_id", Value: bson.D{
						{Key: "$in", Value: filter.Uuids},
					},
				},
			},
		)
	}

	if len(filter.Names) > 0 {
		bsonFilter = append(
			bsonFilter, bson.D{
				{
					Key: "name", Value: bson.D{
						{Key: "$in", Value: filter.Names},
					},
				},
			},
		)
	}

	if len(filter.Categories) > 0 {
		bsonFilter = append(
			bsonFilter, bson.D{
				{
					Key: "category", Value: bson.D{
						{Key: "$in", Value: filter.Categories},
					},
				},
			},
		)
	}

	if len(filter.ManufacturerCountries) > 0 {
		bsonFilter = append(
			bsonFilter, bson.D{
				{
					Key: "manufacturer.country", Value: bson.D{
						{Key: "$in", Value: filter.ManufacturerCountries},
					},
				},
			},
		)
	}

	if len(filter.Tags) > 0 {
		bsonFilter = append(
			bsonFilter, bson.D{
				{
					Key: "tags", Value: bson.D{
						{Key: "$in", Value: filter.Tags},
					},
				},
			},
		)
	}

	if len(bsonFilter) == 0 {
		return bson.D{}
	}

	return bson.D{{Key: "$and", Value: bsonFilter}}
}
