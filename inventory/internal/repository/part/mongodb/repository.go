package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PartStorage struct {
	collection *mongo.Collection
}

func NewPartStorage(mongoDB *mongo.Database) *PartStorage {
	collection := mongoDB.Collection("parts")

	return &PartStorage{
		collection: collection,
	}
}
