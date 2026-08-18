package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/mapper"
	"github.com/paincake00/microservices-go/inventory/internal/repository/model"
)

func (p *PartStorage) GetPart(ctx context.Context, partUUID string) (*entity.Part, error) {
	var part model.Part
	err := p.collection.FindOne(ctx, bson.M{"_id": partUUID}).Decode(&part)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, model.ErrNotFound
	}

	partEntity, err := mapper.FromModel(&part)
	if err != nil {
		return nil, err
	}

	return partEntity, nil
}
