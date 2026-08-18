package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (p *PartStorage) DeleteAll(ctx context.Context) error {
	_, err := p.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}
