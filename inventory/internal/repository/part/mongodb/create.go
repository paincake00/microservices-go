package mongodb

import (
	"context"

	"github.com/paincake00/microservices-go/inventory/internal/entity"
	"github.com/paincake00/microservices-go/inventory/internal/repository/mapper"
)

func (p *PartStorage) Create(ctx context.Context, part *entity.Part) (string, error) {
	model, err := mapper.ToModel(part)
	if err != nil {
		return "", err
	}

	res, err := p.collection.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}
